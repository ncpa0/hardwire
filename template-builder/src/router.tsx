import { ComponentApi, defineContext } from "jsxte";
import path from "node:path";
import { pathCmp, pathIsWithin } from "./utils/paths";

export const appCtx = defineContext<{
  currentRoute: string[];
  selectedRoute: string[];
  addRouter(routerContainerId: string): void;
  registerRoute(path: string, routerContainerId: string): void;
  getRouteContainerId(path: string): string;
}>();

const routerCtx = defineContext<{ containerID: string }>();

export const Router = (
  props: JSXTE.PropsWithChildren<{ id: string }>,
  compApi: ComponentApi
) => {
  const app = compApi.ctx.getOrFail(appCtx);
  app.addRouter(props.id);
  return (
    <div id={props.id}>
      <routerCtx.Provider
        value={{
          containerID: props.id,
        }}
      >
        {props.children}
      </routerCtx.Provider>
    </div>
  );
};

export const Route = (
  props: JSXTE.PropsWithChildren<{
    path: string;
  }>,
  compApi: ComponentApi
) => {
  const app = compApi.ctx.getOrFail(appCtx);
  const router = compApi.ctx.getOrFail(routerCtx);
  app.registerRoute(
    path.join(...app.currentRoute, props.path),
    router.containerID
  );

  if (pathCmp(app.selectedRoute[0] ?? "", props.path)) {
    return (
      <appCtx.Provider
        value={{
          currentRoute: app.currentRoute.concat(props.path),
          selectedRoute: app.selectedRoute.slice(1),
          getRouteContainerId: app.getRouteContainerId,
          registerRoute: app.registerRoute,
          addRouter: app.addRouter,
        }}
      >
        {props.children}
      </appCtx.Provider>
    );
  }

  return <></>;
};

export const Link = (
  props: JSXTE.PropsWithChildren<
    JSX.IntrinsicElements["a"] & {
      href: string;
    }
  >,
  compApi: ComponentApi
) => {
  const app = compApi.ctx.getOrFail(appCtx);
  const currentPathname = app.currentRoute.join("/");

  let base = path.dirname(props.href);
  let name = path.basename(props.href);

  while (true) {
    if (pathIsWithin(currentPathname, base)) {
      const target = app.getRouteContainerId(path.join(base, name));
      return <a {...props} hx-target={target} />;
    }

    let nextBase = path.dirname(base);

    if (base === nextBase || nextBase === "/" || nextBase === "") {
      break;
    }
    name = base.substring(nextBase.length + 1);
    base = nextBase;
  }

  const target = app.getRouteContainerId(currentPathname);
  return <a {...props} hx-target={target} />;
};

export const registerGlobalFunctions = () => {
  Object.defineProperty(global, "Router", {
    value: Router,
    enumerable: true,
    writable: false,
  });
  Object.defineProperty(global, "Route", {
    value: Route,
    enumerable: true,
    writable: false,
  });
  Object.defineProperty(global, "Link", {
    value: Link,
    enumerable: true,
    writable: false,
  });
};
