import { renderToHtmlAsync } from "jsxte";
import prettier from "prettier";
import { collectRoutes } from "./collect-routes";
import { appCtx } from "./router";

const noop = () => {};

export const buildPages = async (tree: JSX.Element) => {
  console.log("Collecting routes...");

  const routes = await collectRoutes(tree);

  const getRouteContainerId = (path: string): string => {
    return routes.get(path)?.containerID ?? routes.topRouter;
  };

  console.log("Building pages...");

  const ops = routes.getAll().map(async (route) => {
    console.log(`Building page '${route.path}.html'`);
    const html = await renderToHtmlAsync(
      <appCtx.Provider
        value={{
          selectedRoute: route.path.split("/"),
          currentRoute: [],
          registerRoute: noop,
          getRouteContainerId,
          addRouter: noop,
        }}
      >
        {tree}
      </appCtx.Provider>
    );

    return {
      route: route.path,
      html: await prettier.format(html, { parser: "html" }),
    };
  });

  return Promise.all(ops);
};
