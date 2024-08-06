import { builderCtx } from "./contexts";
import { render } from "./renderer";
import { capitalize } from "./utils/capitalize";
import { pathCompare } from "./utils/paths";

type RouteDefinition = {
  path: string;
  title: string;
  containerID: string;
};

class RouteCollection {
  private routes: Array<RouteDefinition> = [];
  public topRouter: string = "";

  public add(route: RouteDefinition) {
    this.routes.push(route);
  }

  public get(path: string): RouteDefinition | undefined {
    return this.routes.find((r) => pathCompare(r.path, path));
  }

  public getAll(): Array<RouteDefinition> {
    return this.routes;
  }

  public has(path: string): boolean {
    return this.routes.some((r) => pathCompare(r.path, path));
  }

  public concatInto(collection: RouteCollection) {
    this.routes = this.routes.concat(collection.routes);
  }
}

export const collectRoutes = async (
  entrypointDir: string,
  tree: JSX.Element,
  staticUrl: string,
  selectedRoute: string[] = [],
  collection: RouteCollection = new RouteCollection(),
): Promise<RouteCollection> => {
  const newRoutes: Array<RouteDefinition> = [];

  const registerRoute = (
    path: string,
    title: string,
    routerContainerId: string,
  ) => {
    if (collection.has(path)) {
      return;
    }
    const route = {
      path,
      title,
      containerID: routerContainerId,
    };
    newRoutes.push(route);
    collection.add(route);
  };

  const getRouteContainerId = (path: string): string => {
    return collection.get(path)?.containerID ?? collection.topRouter;
  };

  const addRouter = (routerContainerId: string) => {
    if (!collection.topRouter) {
      collection.topRouter = routerContainerId;
    }
  };

  await render(
    <ExtFilesCtx.Provider value={{ register: () => "", get: () => void 0 }}>
      <builderCtx.Provider
        value={{
          staticUrl: staticUrl,
          isBuildPhase: false,
          entrypointDir,
          selectedRoute,
          currentRoute: [],
          currentRouteTitle: "",
          registerRoute,
          getRouteContainerId,
          addRouter,
          registerDynamicFragment: () => ({ id: "", url: "" }),
          registerRouteDynamicResource: (r) => [capitalize(r), 1],
        }}
      >
        {tree}
      </builderCtx.Provider>
    </ExtFilesCtx.Provider>,
  );

  for (const r of newRoutes) {
    await collectRoutes(
      entrypointDir,
      tree,
      staticUrl,
      r.path.split("/").filter(Boolean),
      collection,
    );
  }

  return collection;
};
