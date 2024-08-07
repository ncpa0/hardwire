import { DynamicFragment } from "./dynamic-fragment";
import { $action } from "./form-action";
import { If } from "./gotmpl-generator/if";
import { MapArray } from "./gotmpl-generator/map-array";
import { Head } from "./head";
import { $island } from "./island";
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
import { $islandList } from "./dynamic-list";
import { Html } from "./html";

export const GLOBALS = {
  Switch,
  StaticRoute,
  DynamicRoute,
  Link,
  Script,
  Style,
  Head,
  Html,
  DynamicFragment,
  If,
  MapArray,
  Redirect,
  Locale,
  Localizations,
  Translate,
  useTranslation,
  $action,
  $island,
  $islandList,
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
      }),
    ),
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
  const Html: typeof GLOBALS.Html;
  const DynamicFragment: typeof GLOBALS.DynamicFragment;
  const If: typeof GLOBALS.If;
  const MapArray: typeof GLOBALS.MapArray;
  const Redirect: typeof GLOBALS.Redirect;
  const Locale: typeof GLOBALS.Locale;
  const Localizations: typeof GLOBALS.Localizations;
  const Translate: typeof GLOBALS.Translate;
  const useTranslation: typeof GLOBALS.useTranslation;
  const $action: typeof GLOBALS.$action;
  const $island: typeof GLOBALS.$island;
  const $islandList: typeof GLOBALS.$islandList;
  type TFunction = TF;
}
