import { ComponentApi } from "jsxte";
import { HtmlContext } from "./html";

export function Head(props: JSX.IntrinsicElements["head"], api: ComponentApi) {
  const html = api.ctx.getOrFail(HtmlContext);

  if (Array.isArray(props.children)) {
    for (const elem of props.children) {
      html.addHeadEntry(<>{elem}</>);
    }
  } else {
    html.addHeadEntry(<>{props.children}</>);
  }

  return <></>;
}
