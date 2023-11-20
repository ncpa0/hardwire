#!/usr/bin/env bun

import fs from "node:fs/promises";
import path from "node:path";
import { Argv } from "./argv";
import { buildPages } from "./build-pages";
import { registerGlobalFunctions } from "./components";

type PageMetadata = {
  isDynamic: boolean;
  resourceName: string;
};

type FragmentMetadata = {
  resourceName: string;
  hash: string;
};

async function main() {
  const argv = new Argv("bldr");

  const build = argv.registerCommand(
    "build",
    {
      src: "string",
      outdir: "string",
      staticdir: "string",
      staticurl: "string",
    },
    async (args) => {
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
        <App />,
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
            resourceName: page.dynamic?.resource ?? "",
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
