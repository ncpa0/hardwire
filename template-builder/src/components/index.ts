import { DynamicFragment } from "./dynamic-fragment";
import { createFormAction } from "./form-action";
import { If } from "./gotmpl-generator/if";
import { MapArray } from "./gotmpl-generator/range";
import { Head } from "./head";
import { Link } from "./link";
import {
  Locale,
  Localizations,
  TFunction as TF,
  Translate,
  useTranslation,
} from "./localizations";
import { Redirect } from "./redirect";
import { DynamicRoute, StaticRoute } from "./route";
import { Switch } from "./router";
import { Script } from "./script";
import { Style } from "./style";

export const GLOBALS = {
  Switch,
  StaticRoute,
  DynamicRoute,
  Link,
  Script,
  Style,
  Head,
  DynamicFragment,
  If,
  MapArray,
  createFormAction,
  Redirect,
  Locale,
  Localizations,
  Translate,
  useTranslation,
};

export const registerGlobalFunctions = () => {
  Object.defineProperties(
    global,
    Object.fromEntries(
      Object.entries(GLOBALS).map(([k, v]) => {
        return [
          k,
          {
            value: v,
            enumerable: true,
            writable: false,
          },
        ];
      })
    )
  );
};

declare global {
  const Switch: typeof GLOBALS.Switch;
  const StaticRoute: typeof GLOBALS.StaticRoute;
  const DynamicRoute: typeof GLOBALS.DynamicRoute;
  const Link: typeof GLOBALS.Link;
  const Script: typeof GLOBALS.Script;
  const Style: typeof GLOBALS.Style;
  const Head: typeof GLOBALS.Head;
  const DynamicFragment: typeof GLOBALS.DynamicFragment;
  const If: typeof GLOBALS.If;
  const MapArray: typeof GLOBALS.MapArray;
  const createFormAction: typeof GLOBALS.createFormAction;
  const Redirect: typeof GLOBALS.Redirect;
  const Locale: typeof GLOBALS.Locale;
  const Localizations: typeof GLOBALS.Localizations;
  const Translate: typeof GLOBALS.Translate;
  const useTranslation: typeof GLOBALS.useTranslation;
  type TFunction = TF;
}
