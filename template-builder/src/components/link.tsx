import { ComponentApi } from "jsxte";
import path from "node:path";
import { builderCtx } from "../contexts";
import { isSubpath } from "../utils/paths";
import { LocaleContext } from "./localizations";

export const Link = (
  props: JSXTE.PropsWithChildren<
    JSX.IntrinsicElements["a"] & {
      href: string;
      /**
       * Overrides the locale of the linked page.
       */
      locale?: string;
    }
  >,
  compApi: ComponentApi
) => {
  let href = props.href;

  const localeCtx = compApi.ctx.get(LocaleContext);
  if (localeCtx && href.startsWith("/")) {
    const lang = props.locale ?? localeCtx.locale;
    href = "/" + path.join(lang, href);
  }

  const app = compApi.ctx.getOrFail(builderCtx);
  const currentPathname = app.currentRoute.join("/");

  let base = path.dirname(href);
  let name = path.basename(href);

  let headers: undefined | string;

  if (props.locale) {
    headers = JSON.stringify({ "Accept-Language": props.locale });
  }

  while (true) {
    if (isSubpath(base, currentPathname)) {
      const target = app.getRouteContainerId(path.join(base, name));
      return (
        <a
          {...props}
          href={href}
          hx-headers={headers}
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
      href={href}
      hx-headers={headers}
      hx-boost="true"
      hx-target={"#" + target}
      hx-swap="outerHTML"
    />
  );
};
