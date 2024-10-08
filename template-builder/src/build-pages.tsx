import path from "node:path";
import { collectRoutes } from "./collect-routes";
import {
  ActionDefinition,
  RegisterExternalFileOptions,
  RouteMetaContext,
  builderCtx,
} from "./contexts";
import { render } from "./renderer";
import { capitalize } from "./utils/capitalize";

const PRETTY_HTML = process.env.PRETTY_HTML;

const noop = () => {};

function createHash(data: string) {
  return Bun.hash.wyhash(data).toString(16);
}

type ExternalFile = {
  contents: string;
  name: string;
  type: string;
  outFile: string;
  url: string;
};

type Page = {
  route: string;
  html: string;
  dynamic?: {
    resources: {
      key: string;
      res: string;
    }[];
  };
  metadata: Record<string, any>;
};

export const buildPages = async (
  entrypointDir: string,
  tree: JSX.Element,
  staticUrl: string,
) => {
  console.log("Collecting routes...");

  const routes = await collectRoutes(entrypointDir, tree, staticUrl);

  const getRouteContainerId = (path: string): string => {
    return routes.get(path)?.containerID ?? routes.topRouter;
  };

  const assets: Array<ExternalFile> = [];

  const registerExternalFile = (
    contents: string,
    name: string,
    type: "js" | "css",
    options?: RegisterExternalFileOptions,
  ) => {
    let hashedName = createHash(`${name}_${contents}`);
    if (options?.keepName) {
      hashedName = `${name}_${hashedName}`;
    }
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
    return { url: `/__dyn/${hashedName}`, id: hashedName };
  };

  const getExternalFileUrl = (name: string) => {
    const file = assets.find((f) => f.name === name);
    return file?.url;
  };

  const actions: Array<ActionDefinition> = [];

  const registerAction = (action: ActionDefinition) => {
    actions.push(action);
  };

  console.log("Building pages...");

  const pages: Array<Page> = [];

  for (const route of routes.getAll()) {
    console.log(`Building page '${route.path}.html'`);

    const page: Page = {
      route: route.path,
      html: "",
      metadata: {},
    };

    let requiredResources: { key: string; res: string }[] = [];

    const registerRouteDynamicResource = (
      resource: string,
    ): [string, number] => {
      const key = capitalize(resource.replace(/[^a-zA-Z]/g, "").trim());
      requiredResources.push({ key, res: resource });
      return [key, requiredResources.length];
    };

    const addMetadata = (key: string, value: any) => {
      page.metadata[key] = value;
    };

    const html = await render(
      <RouteMetaContext.Provider value={{ addMetadata }}>
        <ExtFilesCtx.Provider
          value={{ register: registerExternalFile, get: getExternalFileUrl }}
        >
          <builderCtx.Provider
            value={{
              staticUrl: staticUrl,
              isBuildPhase: true,
              entrypointDir: entrypointDir,
              selectedRoute: route.path.split("/"),
              currentRoute: [],
              currentRouteTitle: route.title,
              registerRoute: noop,
              getRouteContainerId,
              addRouter: noop,
              registerDynamicFragment,
              registerRouteDynamicResource,
              registerAction,
            }}
          >
            {tree}
          </builderCtx.Provider>
        </ExtFilesCtx.Provider>
      </RouteMetaContext.Provider>,
    );

    if (PRETTY_HTML) {
      const prettier = await import("prettier");
      page.html = await prettier.format(html, { parser: "html" });
    } else {
      page.html = html;
    }

    if (requiredResources.length > 0) {
      page.dynamic = {
        resources: requiredResources,
      };
    }

    pages.push(page);
  }

  return { pages, assets, dynamicFragments, actions };
};
