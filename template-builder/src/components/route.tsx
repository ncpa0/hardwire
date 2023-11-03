import { ComponentApi } from "jsxte";
import path from "node:path";
import { builderCtx, routerCtx } from "../contexts";
import { pathCmp } from "../utils/paths";

export const Route = (
  props: JSXTE.PropsWithChildren<{
    path: string;
  }>,
  compApi: ComponentApi
) => {
  const app = compApi.ctx.getOrFail(builderCtx);
  const router = compApi.ctx.getOrFail(routerCtx);
  app.registerRoute(
    path.join(...app.currentRoute, props.path),
    router.containerID
  );

  if (pathCmp(app.selectedRoute[0] ?? "", props.path)) {
    return (
      <builderCtx.Provider
        value={{
          currentRoute: app.currentRoute.concat(props.path),
          selectedRoute: app.selectedRoute.slice(1),
          entrypointDir: app.entrypointDir,
          isBuildPhase: app.isBuildPhase,
          getRouteContainerId: app.getRouteContainerId,
          registerRoute: app.registerRoute,
          addRouter: app.addRouter,
          registerDynamicFragment: app.registerDynamicFragment,
        }}
      >
        {props.children}
      </builderCtx.Provider>
    );
  }

  return <></>;
};
