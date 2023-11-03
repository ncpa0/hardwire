import { ComponentApi } from "jsxte";
import { builderCtx } from "../contexts";
import { StructProxy, structProxy } from "./gotmpl-generator/generate-go-templ";

export type StreamProps<T extends object> = {
  require: string;
  render: (data: StructProxy<T>) => JSX.Element;
};

export const Stream = async <T extends object = Record<never, never>>(
  props: StreamProps<T>,
  compApi: ComponentApi
) => {
  const bldr = compApi.ctx.getOrFail(builderCtx);
  const templ = await compApi.renderAsync(<>{props.render(structProxy(""))}</>);

  const url = bldr.registerDynamicFragment(props.require, templ);

  return <div hx-trigger="load" hx-get={url} />;
};
