import { ComponentApi } from "jsxte";
import path from "node:path/posix";
import { builderCtx, routerCtx } from "../contexts";
import { isSubpath, pathCompare } from "../utils/paths";
import { StructProxy, structProxy } from "./gotmpl-generator/generate-go-templ";

export const StaticRoute = (
  props: JSXTE.PropsWithChildren<{
    path: string;
    exact?: true;
    title?: string;
  }>,
  compApi: ComponentApi
) => {
  const app = compApi.ctx.getOrFail(builderCtx);
  const router = compApi.ctx.getOrFail(routerCtx);
  const thisRoute = path.join(...app.currentRoute, props.path);
  app.registerRoute(thisRoute, props.title ?? "", router.containerID);

  const currentRoute =
    app.selectedRoute.length > 0 ? path.join(...app.selectedRoute) : "";

  if (props.exact) {
    if (!pathCompare(thisRoute, currentRoute)) {
      return <></>;
    }
  }

  if (isSubpath(currentRoute, thisRoute)) {
    return (
      <builderCtx.Provider
        value={{
          currentRoute: app.currentRoute.concat(props.path),
          currentRouteTitle: app.currentRouteTitle,
          selectedRoute: app.selectedRoute,
          entrypointDir: app.entrypointDir,
          isBuildPhase: app.isBuildPhase,
          getRouteContainerId: app.getRouteContainerId,
          registerRoute: app.registerRoute,
          addRouter: app.addRouter,
          registerDynamicFragment: app.registerDynamicFragment,
          registerRouteDynamicResource: app.registerRouteDynamicResource,
        }}
      >
        {props.children}
      </builderCtx.Provider>
    );
  }

  return <></>;
};

export const DynamicRoute = <T extends object = Record<never, never>>(
  props: {
    path: string;
    require: string;
    exact?: boolean;
    title?: string;
    render: (data: StructProxy<T>) => JSX.Element;
  },
  compApi: ComponentApi
) => {
  const app = compApi.ctx.getOrFail(builderCtx);
  const router = compApi.ctx.getOrFail(routerCtx);
  const thisRoute = path.join(...app.currentRoute, props.path);
  app.registerRoute(thisRoute, props.title ?? "", router.containerID);

  const currentRoute =
    app.selectedRoute.length > 0 ? path.join(...app.selectedRoute) : "";

  if (props.exact) {
    if (!pathCompare(thisRoute, currentRoute)) {
      return <></>;
    }
  }

  if (isSubpath(currentRoute, thisRoute)) {
    const [resourceKey, depth] = app.registerRouteDynamicResource(
      props.require
    );
    return (
      <builderCtx.Provider
        value={{
          currentRoute: app.currentRoute.concat(props.path),
          currentRouteTitle: app.currentRouteTitle,
          selectedRoute: app.selectedRoute,
          entrypointDir: app.entrypointDir,
          isBuildPhase: app.isBuildPhase,
          getRouteContainerId: app.getRouteContainerId,
          registerRoute: app.registerRoute,
          addRouter: app.addRouter,
          registerDynamicFragment: app.registerDynamicFragment,
          registerRouteDynamicResource: app.registerRouteDynamicResource,
        }}
      >
        {`{{$root${depth} := .${resourceKey}}}`}
        {props.render(structProxy(`$root${depth}`))}
      </builderCtx.Provider>
    );
  }

  return <></>;
};
