import { ContextDefinition, defineContext } from "jsxte";

export type ActionDefinition = {
  resource: string;
  action: string;
  method: string;
  islandIDs: string[];
};

export const builderCtx = defineContext<{
  staticUrl: string;
  isBuildPhase: boolean;
  entrypointDir: string;
  currentRoute: string[];
  currentRouteTitle: string;
  selectedRoute: string[];
  addRouter(routerContainerId: string): void;
  registerRoute(path: string, title: string, routerContainerId: string): void;
  getRouteContainerId(path: string): string;
  registerDynamicFragment(
    require: string,
    templ: string,
  ): { id: string; url: string };
  registerRouteDynamicResource(require: string): [key: string, deepth: number];
  registerAction(action: ActionDefinition): void;
}>();

export const routerCtx = defineContext<{ containerID: string }>();

export type RegisterExternalFileOptions = {
  keepName?: boolean;
};

declare global {
  const ExtFilesCtx: ContextDefinition<{
    register(
      content: string,
      name: string,
      type: "js" | "css",
      options?: RegisterExternalFileOptions,
    ): string;
    get(name: string): string | undefined;
  }>;
}

export const ExtFilesCtx = defineContext<{
  register(
    content: string,
    name: string,
    type: "js" | "css",
    options?: RegisterExternalFileOptions,
  ): string;
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
