import { ComponentApi } from "jsxte";
import path from "node:path";
import { RouteMetaContext } from "../contexts";
import { LocaleContext } from "./localizations";

export const Redirect = (
  props: {
    href: string;
    /**
     * Overrides the locale of the linked page.
     */
    locale?: string;
  },
  api: ComponentApi
) => {
  const localeCtx = api.ctx.get(LocaleContext);
  const metaCtx = api.ctx.get(RouteMetaContext);

  let href = props.href;
  if (href.startsWith("/")) {
    if (props.locale) {
      href = "/" + path.join(props.locale, href);
    } else if (localeCtx) {
      href = "/" + path.join(localeCtx.locale, href);
    }
  }

  metaCtx?.addMetadata("redirectUrl", href);
  metaCtx?.addMetadata("shouldRedirect", true);

  return <></>;
};
