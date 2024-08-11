import { ComponentApi } from "jsxte";
import { builderCtx } from "../contexts";
import { render } from "../renderer";
import { structProxy } from "./gotmpl-generator/generate-go-templ";
import { cls } from "../utils/cls";
import { IslandDefinition, IslandMap } from "./island";

export type DynamicListProps<T extends any[]> = {
  require: string;
  render: (
    item: AsProxy<T[number]>,
    itemIdx: ValueProxy<string>,
  ) => JSX.Element;
  renderWrapper?: (
    props: JSXTE.PropsWithChildren<{}>,
    items: ListProxy<T[number]>,
  ) => JSX.Element;
  keyGetter?: (
    item: AsProxy<T[number]>,
    idx: ValueProxy<string>,
  ) => string | ValueProxy<string>;
  id?: string;
  class?: string | ValueProxy<string>;
  itemClass?: string | ValueProxy<string>;
  fallback?: JSX.Element;
  trigger?: "load" | "revealed" | "intersect";
  morph?: boolean;
  swap?: `${number}s` | `${number}ms`;
  settle?: `${number}s` | `${number}ms`;
  /**
   * Overrides the accepted language header of the request.
   */
  locale?: string;
};

export const DynamicList = async <T extends any[]>(
  props: DynamicListProps<T>,
  compApi: ComponentApi,
) => {
  let itemClassName = "dynamic-list-element";
  if (props.itemClass) {
    itemClassName += " " + props.itemClass;
  }
  const getKey = props.keyGetter || ((_, idx) => idx);

  const bldr = compApi.ctx.getOrFail(builderCtx);
  const listProxy = structProxy<T>("$frag_root") as ListProxy<T>;
  const itemsTempl = listProxy.map((value, key) => {
    const itemKey = getKey(value, key);
    return (
      <div
        data-item-key={itemKey}
        class={cls("dynamic-list-element", props.itemClass)}
      >
        {props.render(value, key)}
      </div>
    );
  });
  const templ = await render(
    <dynamic-fragment id={props.id} class={cls("dynamic-list", props.class)}>
      {`{{$frag_root := .}}`}
      {props.renderWrapper
        ? props.renderWrapper({ children: itemsTempl }, listProxy)
        : itemsTempl}
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

  let hxswap = props.morph ? "morph:outerHTML" : "outerHTML";

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

export function $islandList<T extends any[], P extends object = {}>(
  options: Omit<DynamicListProps<T>, "render"> & { id: string },
  ItemComponent: (
    props: P,
    item: AsProxy<T[number]>,
    itemIdx: ValueProxy<string>,
  ) => JSX.Element,
  WrapperComponent?: (
    props: JSXTE.PropsWithChildren<{}>,
    items: ListProxy<T[number]>,
  ) => JSX.Element,
): JSXTE.Component<P> {
  const { id, ...dynamicFragmentProps } = options;

  const islandEntry: IslandDefinition = {
    id,
    fragmentID: "",
    type: "list",
    resource: dynamicFragmentProps.require,
  };

  const island: JSXTE.Component<P> = (props) => {
    return (
      <DynamicList<T>
        id={id}
        class={`island_${id}`}
        {...dynamicFragmentProps}
        render={(item, itemIdx) => {
          return ItemComponent(props, item, itemIdx);
        }}
        renderWrapper={WrapperComponent}
        // @ts-expect-error
        __fragidgetter={(id: string) => {
          islandEntry.fragmentID = id;
        }}
      />
    );
  };

  IslandMap.set(island, islandEntry);

  return island;
}
