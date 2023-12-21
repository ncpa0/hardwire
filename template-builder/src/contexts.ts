import { ContextDefinition, defineContext } from "jsxte";

export const builderCtx = defineContext<{
  isBuildPhase: boolean;
  entrypointDir: string;
  currentRoute: string[];
  currentRouteTitle: string;
  selectedRoute: string[];
  addRouter(routerContainerId: string): void;
  registerRoute(path: string, title: string, routerContainerId: string): void;
  getRouteContainerId(path: string): string;
  registerDynamicFragment(require: string, templ: string): string;
  registerRouteDynamicResource(require: string): [key: string, deepth: number];
}>();

export const routerCtx = defineContext<{ containerID: string }>();

declare global {
  const ExtFilesCtx: ContextDefinition<{
    register(content: string, name: string, type: string): string;
    get(name: string): string | undefined;
  }>;
}

export const ExtFilesCtx = defineContext<{
  register(content: string, name: string, type: string): string;
  get(name: string): string;
}>();

Object.defineProperty(global, "ExtFilesCtx", {
  value: ExtFilesCtx,
  writable: false,
  enumerable: false,
  configurable: false,
});

export const RouteMetaContext = defineContext<{
  addMetadata(key: string, value: any): void;
}>();
