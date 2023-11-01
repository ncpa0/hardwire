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

  const extFiles: Array<ExternalFile> = [];

  const registerExternalFile = (
    contents: string,
    name: string,
    type: "js" | "css"
  ) => {
    const hashedName = createHash(`${name}_${contents}`);
    switch (type) {
      case "css": {
        const fPath = `/assets/css/${hashedName}.css`;
        extFiles.push({
          contents,
          name,
          type,
          outFile: fPath,
        });
        return path.join(staticUrl, fPath);
      }
      case "js": {
        const fPath = `/assets/js/${hashedName}.js`;
        extFiles.push({
          contents,
          name,
          type,
          outFile: fPath,
        });
        return path.join(staticUrl, fPath);
      }
    }
  };

  console.log("Building pages...");

  const ops = routes.getAll().map(async (route) => {
    console.log(`Building page '${route.path}.html'`);
    const html = await renderToHtmlAsync(
      <ExtFilesCtx.Provider value={{ register: registerExternalFile }}>
        <builderCtx.Provider
          value={{
            isBuildStep: true,
            entrypointDir: entrypointDir,
            selectedRoute: route.path.split("/"),
            currentRoute: [],
            registerRoute: noop,
            getRouteContainerId,
            addRouter: noop,
          }}
        >
          {tree}
        </builderCtx.Provider>
      </ExtFilesCtx.Provider>
    );

    return {
      route: route.path,
      html: await prettier.format(html, { parser: "html" }),
    };
  });

  return { htmlFiles: await Promise.all(ops), assets: extFiles };
};
