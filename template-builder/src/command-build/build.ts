import { createElement } from "jsxte";
import fs from "node:fs/promises";
import path from "node:path";

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

export async function buildCmd(
  srcFile: string,
  outDir: string,
  staticDir: string,
  staticurl: string
) {
  const { buildPages } = await import("../build-pages");
  const { registerGlobalFunctions } = await import("../components");

  registerGlobalFunctions();

  const mod = await import("file://" + srcFile);

  if (typeof mod.default !== "function") {
    throw new Error(`Module ${srcFile} does not export a default function.`);
  }

  const App = mod.default;

  const { pages, dynamicFragments, assets } = await buildPages(
    path.dirname(srcFile),
    createElement(App),
    staticurl
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
