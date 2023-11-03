import { If } from "./gotmpl-generator/if";
import { Range } from "./gotmpl-generator/range";
import { Head } from "./head";
import { Link } from "./link";
import { Route } from "./route";
import { Router } from "./router";
import { Script } from "./script";
import { Stream } from "./stream";
import { Style } from "./style";

export const registerGlobalFunctions = () => {
  Object.defineProperty(global, "Router", {
    value: Router,
    enumerable: true,
    writable: false,
  });
  Object.defineProperty(global, "Route", {
    value: Route,
    enumerable: true,
    writable: false,
  });
  Object.defineProperty(global, "Link", {
    value: Link,
    enumerable: true,
    writable: false,
  });
  Object.defineProperty(global, "Script", {
    value: Script,
    enumerable: true,
    writable: false,
  });
  Object.defineProperty(global, "Style", {
    value: Style,
    enumerable: true,
    writable: false,
  });
  Object.defineProperty(global, "Head", {
    value: Head,
    enumerable: true,
    writable: false,
  });

  Object.defineProperty(global, "Stream", {
    value: Stream,
    enumerable: true,
    writable: false,
  });
  Object.defineProperty(global, "If", {
    value: If,
    enumerable: true,
    writable: false,
  });
  Object.defineProperty(global, "Range", {
    value: Range,
    enumerable: true,
    writable: false,
  });
};
