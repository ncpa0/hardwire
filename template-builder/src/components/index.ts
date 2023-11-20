import { createFormAction } from "./form-action";
import { If } from "./gotmpl-generator/if";
import { MapArray } from "./gotmpl-generator/range";
import { Head } from "./head";
import { Link } from "./link";
import { DynamicRoute, StaticRoute } from "./route";
import { Switch } from "./router";
import { Script } from "./script";
import { DynamicFragment } from "./stream";
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
}
