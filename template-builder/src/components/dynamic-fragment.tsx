import { ComponentApi } from "jsxte";
import { builderCtx } from "../contexts";
import { render } from "../renderer";
import { structProxy } from "./gotmpl-generator/generate-go-templ";

declare global {
  namespace JSX {
    interface IntrinsicElements {
      "dynamic-fragment": {
        class?: string | ValueProxy<string>;
        children?: JSX.Element;
      };
    }
  }
}

export type DynamicFragmentProps<T> = {
  require: string;
  render: (data: AsProxy<T>) => JSX.Element;
  class?: string;
  fallback?: JSX.Element;
  trigger?: "load" | "revealed" | "intersect";
  swap?: `${number}s` | `${number}ms`;
  settle?: `${number}s` | `${number}ms`;
  /**
   * Overrides the accepted language header of the request.
   */
  locale?: string;
};

export const DynamicFragment = async <T,>(
  props: DynamicFragmentProps<T>,
  compApi: ComponentApi,
) => {
  const bldr = compApi.ctx.getOrFail(builderCtx);
  const templ = await render(
    <dynamic-fragment class={props.class}>
      {`{{$frag_root := .}}`}
      {props.render(structProxy("$frag_root"))}
    </dynamic-fragment>,
    compApi,
  );

  const { url, id } = bldr.registerDynamicFragment(props.require, templ);

  if ("__fragidgetter" in props) {
    (props as any).__fragidgetter(id);
  }

  let hxtrigger = "revealed";
  switch (props.trigger) {
    case "load":
      hxtrigger = "load delay:20ms";
      break;
    case "revealed":
      hxtrigger = "revealed delay:20ms";
      break;
    case "intersect":
      hxtrigger = "intersect delay:20ms";
      break;
  }

  let hxswap = "outerHTML";

  if (props.swap) {
    hxswap += " swap:" + props.swap;
  }

  if (props.settle) {
    hxswap += " settle:" + props.settle;
  }

  const headers: Record<string, string> = {
    "Hardwire-Dynamic-Fragment-Request": "/" + bldr.currentRoute.join("/"),
  };

  if (props.locale) {
    headers["Accept-Language"] = props.locale;
  }

  return (
    <div
      hx-trigger={hxtrigger}
      hx-get={url}
      hx-swap={hxswap}
      hx-headers={JSON.stringify(headers)}
    >
      {props.fallback}
    </div>
  );
};
