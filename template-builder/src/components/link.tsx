import { ComponentApi } from "jsxte";
import path from "node:path";
import { builderCtx } from "../contexts";
import { pathIsWithin } from "../utils/paths";

export const Link = (
  props: JSXTE.PropsWithChildren<
    JSX.IntrinsicElements["a"] & {
      href: string;
    }
  >,
  compApi: ComponentApi
) => {
  const app = compApi.ctx.getOrFail(builderCtx);
  const currentPathname = app.currentRoute.join("/");

  let base = path.dirname(props.href);
  let name = path.basename(props.href);

  while (true) {
    if (pathIsWithin(currentPathname, base)) {
      const target = app.getRouteContainerId(path.join(base, name));
      return (
        <a
          {...props}
          hx-boost="true"
          hx-target={"#" + target}
          hx-swap="outerHTML"
        />
      );
    }

    let nextBase = path.dirname(base);

    if (base === nextBase || nextBase === "/" || nextBase === "") {
      break;
    }
    name = base.substring(nextBase.length + 1);
    base = nextBase;
  }

  const target = app.getRouteContainerId(currentPathname);
  return (
    <a
      {...props}
      hx-boost="true"
      hx-target={"#" + target}
      hx-swap="outerHTML"
    />
  );
};
