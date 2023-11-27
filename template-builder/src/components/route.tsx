import { ComponentApi } from "jsxte";
import path from "node:path";
import { builderCtx, routerCtx } from "../contexts";
import { pathCmp } from "../utils/paths";
import { StructProxy, structProxy } from "./gotmpl-generator/generate-go-templ";

export const StaticRoute = (
  props: JSXTE.PropsWithChildren<{
    path: string;
    title?: string;
  }>,
  compApi: ComponentApi
) => {
  const app = compApi.ctx.getOrFail(builderCtx);
  const router = compApi.ctx.getOrFail(routerCtx);
  app.registerRoute(
    path.join(...app.currentRoute, props.path),
    props.title ?? "",
    router.containerID
  );

  if (pathCmp(app.selectedRoute[0] ?? "", props.path)) {
    return (
      <builderCtx.Provider
        value={{
          currentRoute: app.currentRoute.concat(props.path),
          currentRouteTitle: app.currentRouteTitle,
          selectedRoute: app.selectedRoute.slice(1),
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
    title?: string;
    render: (data: StructProxy<T>) => JSX.Element;
  },
  compApi: ComponentApi
) => {
  const app = compApi.ctx.getOrFail(builderCtx);
  const router = compApi.ctx.getOrFail(routerCtx);
  app.registerRoute(
    path.join(...app.currentRoute, props.path),
    props.title ?? "",
    router.containerID
  );

  if (pathCmp(app.selectedRoute[0] ?? "", props.path)) {
    const [resourceKey, depth] = app.registerRouteDynamicResource(
      props.require
    );
    return (
      <builderCtx.Provider
        value={{
          currentRoute: app.currentRoute.concat(props.path),
          currentRouteTitle: app.currentRouteTitle,
          selectedRoute: app.selectedRoute.slice(1),
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
