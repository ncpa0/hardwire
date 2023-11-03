import { renderToHtmlAsync } from "jsxte";
import path from "node:path";
import prettier from "prettier";
import { collectRoutes } from "./collect-routes";
import { builderCtx } from "./contexts";

const noop = () => {};

function createHash(data: string) {
  return Bun.hash.crc32(data).toString(16);
}

type ExternalFile = {
  contents: string;
  name: string;
  type: string;
  outFile: string;
  url: string;
};

export const buildPages = async (
  entrypointDir: string,
  tree: JSX.Element,
  staticUrl: string
) => {
  console.log("Collecting routes...");

  const routes = await collectRoutes(entrypointDir, tree);

  const getRouteContainerId = (path: string): string => {
    return routes.get(path)?.containerID ?? routes.topRouter;
  };

  const assets: Array<ExternalFile> = [];

  const registerExternalFile = (
    contents: string,
    name: string,
    type: "js" | "css"
  ) => {
    const hashedName = createHash(`${name}_${contents}`);
    switch (type) {
      case "css": {
        const fPath = `/assets/css/${hashedName}.css`;
        const url = path.join(staticUrl, fPath);
        assets.push({
          contents,
          name,
          type,
          outFile: fPath,
          url,
        });
        return url;
      }
      case "js": {
        const fPath = `/assets/js/${hashedName}.js`;
        const url = path.join(staticUrl, fPath);
        assets.push({
          contents,
          name,
          type,
          outFile: fPath,
          url,
        });
        return url;
      }
    }
  };

  const dynamicFragments: Array<{
    name: string;
    hash: string;
    contents: string;
  }> = [];

  const registerDynamicFragment = (name: string, contents: string) => {
    const hashedName = createHash(`${name}_${contents}`);
    dynamicFragments.push({
      name,
      hash: hashedName,
      contents,
    });
    return `/__htmxdyn_frag/${hashedName}`;
  };

  const getExternalFileUrl = (name: string) => {
    const file = assets.find((f) => f.name === name);
    return file?.url;
  };

  console.log("Building pages...");

  const htmlFiles: Array<{ route: string; html: string }> = [];

  for (const route of routes.getAll()) {
    console.log(`Building page '${route.path}.html'`);
    const html = await renderToHtmlAsync(
      <ExtFilesCtx.Provider
        value={{ register: registerExternalFile, get: getExternalFileUrl }}
      >
        <builderCtx.Provider
          value={{
            isBuildPhase: true,
            entrypointDir: entrypointDir,
            selectedRoute: route.path.split("/"),
            currentRoute: [],
            registerRoute: noop,
            getRouteContainerId,
            addRouter: noop,
            registerDynamicFragment,
          }}
        >
          {tree}
        </builderCtx.Provider>
      </ExtFilesCtx.Provider>
    );

    htmlFiles.push({
      route: route.path,
      html: await prettier.format(html, { parser: "html" }),
    });
  }

  return { htmlFiles, assets, dynamicFragments };
};
