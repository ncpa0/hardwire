import { If } from "./gotmpl-generator/if";
import { RangeOver } from "./gotmpl-generator/range";
import { Head } from "./head";
import { Link } from "./link";
import { Route } from "./route";
import { Router } from "./router";
import { Script } from "./script";
import { Stream } from "./stream";
import { Style } from "./style";

export const GLOBALS = {
  Router,
  Route,
  Link,
  Script,
  Style,
  Head,
  Stream,
  If,
  RangeOver,
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
  const Router: typeof GLOBALS.Router;
  const Route: typeof GLOBALS.Route;
  const Link: typeof GLOBALS.Link;
  const Script: typeof GLOBALS.Script;
  const Style: typeof GLOBALS.Style;
  const Head: typeof GLOBALS.Head;
  const Stream: typeof GLOBALS.Stream;
  const If: typeof GLOBALS.If;
  const RangeOver: typeof GLOBALS.RangeOver;
}
