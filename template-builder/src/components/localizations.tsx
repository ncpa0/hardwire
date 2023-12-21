import { ComponentApi, defineContext } from "jsxte";
import { StaticRoute } from "./route";
import { Switch } from "./router";

export type TFunction = (
  key: string,
  params?: Record<string, string | number>
) => string;

export const LocaleContext = defineContext<{
  locale: string;
  tFunction: TFunction;
}>();

export const Localizations: JSXTE.Component = (props) => {
  return <Switch id="__localization_route_switch">{props.children}</Switch>;
};

export type LocaleProps = {
  lang: string;
  translations: Record<string, any>;
};

const tFunctionFactory = (translations: Record<string, any>): TFunction => {
  const interpolateParams = (
    str: string,
    params: Record<string, string | number>
  ) => {
    const strSegments: Array<
      { type: "string"; value: string } | { type: "tkey"; value: string }
    > = [];
    let s = str;
    let nextParamIdx: number;
    while ((nextParamIdx = s.search(/{{.*?}}/)) != -1) {
      const closingIdx = s.indexOf("}}");
      strSegments.push({ type: "string", value: s.substring(0, nextParamIdx) });
      strSegments.push({
        type: "tkey",
        value: s.substring(nextParamIdx + 2, closingIdx),
      });
      s = s.substring(closingIdx + 2);
    }

    return (
      strSegments
        .map((seg) => {
          if (seg.type === "string") {
            return seg.value;
          }
          return params[seg.value] ?? seg.value;
        })
        .join("") + s
    );
  };

  const tFunction: TFunction = (key, params) => {
    const pathSegments = key.split(".");
    let translation = translations;
    for (const segment of pathSegments) {
      translation = translation[segment];
    }
    if (typeof translation !== "string") {
      console.warn(`Translation not found for key: ${key}`);
      return key;
    }
    if (params) {
      return interpolateParams(translation, params);
    }
    return translation as string;
  };

  return tFunction;
};

export const Locale: JSXTE.Component<LocaleProps> = (props) => {
  const tFunction = tFunctionFactory(props.translations);

  return (
    <StaticRoute path={props.lang}>
      <LocaleContext.Provider
        value={{
          locale: props.lang,
          tFunction,
        }}
      >
        <div
          class={`__hardwire_locale locale_${props.lang}`}
          hx-headers={JSON.stringify({ "Accept-Language": props.lang })}
        >
          {props.children}
        </div>
      </LocaleContext.Provider>
    </StaticRoute>
  );
};

export const Translate: JSXTE.Component<{
  params?: Record<string, string | number>;
}> = (props, api) => {
  const ctx = api.ctx.get(LocaleContext);

  const t = ctx?.tFunction ?? ((key: string) => key);
  const tc = (elem: JSX.Element | JSX.Element[]) => {
    if (typeof elem === "string") {
      const keys = elem.split(" ").map((k) => k.trim());
      return keys.map((k) => t(k, props.params)).join(" ");
    }
    if (
      typeof elem === "object" &&
      elem != null &&
      !Array.isArray(elem) &&
      (elem as JSXTE.TagElement | JSXTE.TextNodeElement).type === "textNode"
    ) {
      (elem as JSXTE.TextNodeElement).text = t(
        (elem as JSXTE.TextNodeElement).text,
        props.params
      );
    }
    return elem;
  };

  const children = Array.isArray(props.children)
    ? props.children.map(tc)
    : tc(props.children);

  return <>{children}</>;
};

export const useTranslation = (componentApi: ComponentApi): TFunction => {
  const ctx = componentApi.ctx.get(LocaleContext);
  const t = ctx?.tFunction ?? ((key: string) => key);
  return t;
};
