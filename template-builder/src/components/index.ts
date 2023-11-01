import { Link } from "./link";
import { Route } from "./route";
import { Router } from "./router";
import { Script } from "./script";
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
};
