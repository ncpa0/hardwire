#!/usr/bin/env bun

import { createElement } from "jsxte";
import fs from "node:fs/promises";
import path from "node:path";
import { Argv } from "./argv";

type PageMetadata = {
  isDynamic: boolean;
  resources?: {
    key: string;
    res: string;
  }[];
};

type FragmentMetadata = {
  resourceName: string;
  hash: string;
};

async function main() {
  const argv = new Argv("hardwire-html-builder");

  const init = argv.registerCommand(
    "init",
    {
      dir: "string",
    },
    async (args) => {
      const wd = path.resolve(args.dir);

      const bunfigPath = path.join(wd, "bunfig.toml");
      if (!(await fs.exists(bunfigPath))) {
        const bunfig = `jsx = "reacjsxt-jsx"\njsxImportSource = "jsxte"\n`;
        await Bun.write(bunfigPath, bunfig);
      }

      const indexPath = path.join(wd, "index.tsx");
      if (!(await fs.exists(indexPath))) {
        const content = `export default function App() {
  return (
    <html>
      <Head>
        <meta charset='utf-8' />
        <meta http-equiv='X-UA-Compatible' content='IE=edge' />
        <meta name='viewport' content='width=device-width, initial-scale=1' />
        <Style path="./style.css" dirname={__dirname} />
        <Script path="./index.client.ts" dirname={__dirname} />
      </Head>
      <body>
        <div id="root">
          <nav>
            <ul>
              <li><Link href="/">Home</Link></li>
            </ul>
          </nav>
          <Switch>
            <Route path="home" title="Home Page">
              <h1>Hello World</h1>
            </Route>
          </Switch>
        </div>
      </body>
    </html>
  );
}`;

        await Bun.write(indexPath, content);

        const stylePath = path.join(wd, "style.css");
        if (!(await fs.exists(stylePath))) {
          const content = `body { 
  margin: unset;
  min-height: 100vh;
  min-width: 100vw;
}

#root {
  display: flex;
  flex-direction: column;
}`;
          await Bun.write(stylePath, content);
        }

        const indexClientPath = path.join(wd, "index.client.ts");
        if (!(await fs.exists(indexClientPath))) {
          await Bun.write(indexClientPath, "");
        }
      }
    }
  );

  const build = argv.registerCommand(
    "build",
    {
      src: "string",
      outdir: "string",
      staticdir: "string",
      staticurl: "string",
    },
    async (args) => {
      const { buildPages } = await import("./build-pages");
      const { registerGlobalFunctions } = await import("./components");

      const srcFile = path.resolve(args.src);
      const outDir = path.resolve(args.outdir);
      const staticDir = path.resolve(args.staticdir);

      registerGlobalFunctions();

      const mod = await import("file://" + srcFile);

      if (typeof mod.default !== "function") {
        throw new Error(
          `Module ${srcFile} does not export a default function.`
        );
      }

      const App = mod.default;

      const { pages, dynamicFragments, assets } = await buildPages(
        path.dirname(srcFile),
        createElement(App),
        args.staticurl
      );

      console.log("Saving results to filesystem...");
      await Promise.all([
        ...pages.map(async (page) => {
          const outfilePath = path.join(outDir, page.route) + ".html";
          const metaFilePath = path.join(outDir, page.route) + ".meta.json";
          const basedir = path.dirname(outfilePath);
          const meta: PageMetadata = {
            isDynamic: page.dynamic != null,
            resources: page.dynamic?.resources,
          };

          await fs.mkdir(basedir, { recursive: true });
          await Bun.write(outfilePath, page.html);
          await Bun.write(metaFilePath, JSON.stringify(meta));
        }),
        ...assets.map(async (asset) => {
          const outfilePath = path.join(staticDir, asset.outFile);
          const basedir = path.dirname(outfilePath);
          await fs.mkdir(basedir, { recursive: true });
          await Bun.write(outfilePath, asset.contents);
        }),
        ...dynamicFragments.map(async (frag) => {
          const meta: FragmentMetadata = {
            resourceName: frag.name,
            hash: frag.hash,
          };

          const outfilePath =
            path.join(outDir, "__dyn", frag.hash) + ".template.html";
          const metaFile = path.join(outDir, "__dyn", frag.hash) + ".meta.json";

          const basedir = path.dirname(outfilePath);
          await fs.mkdir(basedir, { recursive: true });

          await Bun.write(outfilePath, frag.contents);
          await Bun.write(metaFile, JSON.stringify(meta));
        }),
      ]);

      console.log("Done.");
    }
  );

  build.setDescription(
    "Generates static HTML files from the given JSX component."
  );
  argv.run();
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
