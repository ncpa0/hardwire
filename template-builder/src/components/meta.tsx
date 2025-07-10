import { ComponentApi } from "jsxte";
import { HtmlContext } from "./html";

export function Meta(props: JSX.IntrinsicElements["meta"], api: ComponentApi) {
  const html = api.ctx.getOrFail(HtmlContext);
  html.addHeadEntry(<meta {...props} />);
  return <></>;
}
