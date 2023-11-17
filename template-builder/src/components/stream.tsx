import { ComponentApi } from "jsxte";
import { builderCtx } from "../contexts";
import { StructProxy, structProxy } from "./gotmpl-generator/generate-go-templ";

declare global {
  namespace JSX {
    interface IntrinsicElements {
      "dynamic-fragment": {
        class?: string;
        children?: JSX.Element;
      };
    }
  }
}

export type DynamicFragmentProps<T extends object> = {
  require: string;
  render: (data: StructProxy<T>) => JSX.Element;
  class?: string;
  fallback?: JSX.Element;
  trigger?: "load" | "revealed";
  swap?: `${number}s` | `${number}ms`;
  settle?: `${number}s` | `${number}ms`;
};

export const DynamicFragment = async <T extends object = Record<never, never>>(
  props: DynamicFragmentProps<T>,
  compApi: ComponentApi
) => {
  const bldr = compApi.ctx.getOrFail(builderCtx);
  const templ = await compApi.renderAsync(
    <dynamic-fragment class={props.class}>
      {props.render(structProxy(""))}
    </dynamic-fragment>
  );

  const url = bldr.registerDynamicFragment(props.require, templ);

  const hxtrigger =
    props.trigger === "load" ? "load delay:20ms" : "revealed delay:20ms";

  let hxswap = "outerHTML";

  if (props.swap) {
    hxswap += " swap:" + props.swap;
  }

  if (props.settle) {
    hxswap += " settle:" + props.settle;
  }

  return (
    <div hx-trigger={hxtrigger} hx-get={url} hx-swap={hxswap}>
      {props.fallback}
    </div>
  );
};
