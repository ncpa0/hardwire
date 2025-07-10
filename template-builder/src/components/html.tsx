import { builderCtx } from "../contexts";
import { SetupClientHelpers } from "./_helpers.client";
import { ComponentApi, defineContext } from "jsxte";
import { PkgOpt } from "./script";

export interface HtmxConfig {
  /**
   * The attributes to settle during the settling phase.
   * @default ["class", "style", "width", "height"]
   */
  attributesToSettle?: ["class", "style", "width", "height"] | string[];
  /**
   * If the focused element should be scrolled into view.
   * @default false
   */
  defaultFocusScroll?: boolean;
  /**
   * The default delay between completing the content swap and settling attributes.
   * @default 20
   */
  defaultSettleDelay?: number;
  /**
   * The default delay between receiving a response from the server and doing the swap.
   * @default 0
   */
  defaultSwapDelay?: number;
  /**
   * The default swap style to use if **[hx-swap](https://htmx.org/attributes/hx-swap)** is omitted.
   * @default "innerHTML"
   */
  defaultSwapStyle?: "innerHTML" | string;
  /**
   * The number of pages to keep in **localStorage** for history support.
   * @default 10
   */
  historyCacheSize?: number;
  /**
   * Whether or not to use history.
   * @default true
   */
  historyEnabled?: boolean;
  /**
   * If true, htmx will inject a small amount of CSS into the page to make indicators invisible unless the **htmx-indicator** class is present.
   * @default true
   */
  includeIndicatorStyles?: boolean;
  /**
   * The class to place on indicators when a request is in flight.
   * @default "htmx-indicator"
   */
  indicatorClass?: "htmx-indicator" | string;
  /**
   * The class to place on triggering elements when a request is in flight.
   * @default "htmx-request"
   */
  requestClass?: "htmx-request" | string;
  /**
   * The class to temporarily place on elements that htmx has added to the DOM.
   * @default "htmx-added"
   */
  addedClass?: "htmx-added" | string;
  /**
   * The class to place on target elements when htmx is in the settling phase.
   * @default "htmx-settling"
   */
  settlingClass?: "htmx-settling" | string;
  /**
   * The class to place on target elements when htmx is in the swapping phase.
   * @default "htmx-swapping"
   */
  swappingClass?: "htmx-swapping" | string;
  /**
   * Allows the use of eval-like functionality in htmx, to enable **hx-vars**, trigger conditions & script tag evaluation. Can be set to **false** for CSP compatibility.
   * @default true
   */
  allowEval?: boolean;
  /**
   * Use HTML template tags for parsing content from the server. This allows you to use Out of Band content when returning things like table rows, but it is *not* IE11 compatible.
   * @default false
   */
  useTemplateFragments?: boolean;
  /**
   * Allow cross-site Access-Control requests using credentials such as cookies, authorization headers or TLS client certificates.
   * @default false
   */
  withCredentials?: boolean;
  /**
   * The default implementation of **getWebSocketReconnectDelay** for reconnecting after unexpected connection loss by the event code **Abnormal Closure**, **Service Restart** or **Try Again Later**.
   * @default "full-jitter"
   */
  wsReconnectDelay?: "full-jitter" | string | ((retryCount: number) => number);
  // following don't appear in the docs
  /** @default false */
  refreshOnHistoryMiss?: boolean;
  /** @default 0 */
  timeout?: number;
  /** @default "[hx-disable], [data-hx-disable]" */
  disableSelector?: "[hx-disable], [data-hx-disable]" | string;
  /** @default "smooth" */
  scrollBehavior?: "smooth" | "auto";
  /**
   * If set to false, disables the interpretation of script tags.
   * @default true
   */
  allowScriptTags?: boolean;
  /**
   * If set to true, disables htmx-based requests to non-origin hosts.
   * @default false
   */
  selfRequestsOnly?: boolean;
  /**
   * Whether or not the target of a boosted element is scrolled into the viewport.
   * @default true
   */
  scrollIntoViewOnBoost?: boolean;
  /**
   * If set, the nonce will be added to inline scripts.
   * @default ''
   */
  inlineScriptNonce?: string;
  /**
   * The type of binary data being received over the WebSocket connection
   * @default 'blob'
   */
  wsBinaryType?: "blob" | "arraybuffer";
  /**
   * If set to true htmx will include a cache-busting parameter in GET requests to avoid caching partial responses by the browser
   * @default false
   */
  getCacheBusterParam?: boolean;
  /**
   * If set to true, htmx will use the View Transition API when swapping in new content.
   * @default false
   */
  globalViewTransitions?: boolean;
  /**
   * htmx will format requests with these methods by encoding their parameters in the URL, not the request body
   * @default ["get"]
   */
  methodsThatUseUrlParams?: (
    | "get"
    | "head"
    | "post"
    | "put"
    | "delete"
    | "connect"
    | "options"
    | "trace"
    | "patch"
  )[];
  /**
   * If set to true htmx will not update the title of the document when a title tag is found in new content
   * @default false
   */
  ignoreTitle?: boolean;
}

export type HtmlProps = JSXTE.PropsWithChildren<{
  headOptions?: Omit<InternalHeadProps, "extensions">;
  htmlProps?: JSX.IntrinsicElements["html"];
  /**
   * Don't include the idiomorph htmx extension.
   */
  nomorph?: boolean;
  /**
   * Don't include the chunked-transfer htmx extension.
   */
  nochunked?: boolean;
  /**
   * Don't include the htmx-ext-head-support htmx extension.
   */
  nohead?: boolean;
}>;

export const HtmlContext = defineContext<{
  addHeadEntry(head: JSX.Element): void;
}>();

export async function Html(
  {
    children,
    headOptions: headProps,
    htmlProps,
    nochunked,
    nomorph,
    nohead,
  }: HtmlProps,
  api: ComponentApi,
) {
  let extensions: string[] = [];
  if (nomorph !== true) {
    extensions.push("morph");
  }
  if (nochunked !== true) {
    extensions.push("chunked-transfer");
  }
  if (nohead !== true) {
    extensions.push("head-support");
  }

  const headContent: JSX.Element[] = [];

  const contents = await api.renderAsync(
    <HtmlContext.Provider
      value={{
        addHeadEntry(elem) {
          headContent.push(elem);
        },
      }}
    >
      {children}
    </HtmlContext.Provider>,
  );

  return (
    <html {...htmlProps}>
      <InternalHead
        {...headProps}
        extensions={extensions}
        htmxConfig={{ includeIndicatorStyles: false, ...headProps?.htmxConfig }}
      >
        {headContent}
      </InternalHead>
      <body hx-ext={extensions.join(",")}>{contents}</body>
    </html>
  );
}

type InternalHeadProps = {
  /**
   * The integrity hash to use for the htmx script. Versions between 1.9.0
   * and 1.9.10 have their hashes hardcoded in, so they don't need this
   * value provided.
   */
  htmxIntegrityHash?: string;
  extensions?: string[];
  htmxConfig: HtmxConfig;
};

const InternalHead: JSXTE.Component<InternalHeadProps> = (
  props,
  componentApi,
) => {
  const app = componentApi.ctx.getOrFail(builderCtx);
  const extFiles = componentApi.ctx.getOrFail(ExtFilesCtx);

  const hwScript = extFiles.register(
    `(${String(SetupClientHelpers)})(${JSON.stringify(props.extensions ?? [])})`,
    "hardwire",
    "js",
    { keepName: true },
  );

  const packages: Array<string | PkgOpt> = [];
  props.extensions?.forEach((ext) => {
    switch (ext) {
      case "chunked-transfer": {
        packages.push("htmx.ext...chunked-transfer");
        break;
      }
      case "morph": {
        packages.push("idiomorph");
        break;
      }
      case "head-support": {
        packages.push("htmx-ext-head-support/head-support.js");
        break;
      }
    }
  });

  return (
    <head>
      <meta charset="utf-8" />
      <meta http-equiv="x-ua-compatible" content="IE=edge" />
      <meta name="viewport" content="width=device-width, initial-scale=1" />
      <meta name="htmx-config" content={JSON.stringify(props.htmxConfig)} />
      <Script package={{ name: "htmx.org", global: "htmx" }} />
      <Script package={packages} />
      <script src={hwScript} />
      <>{props.children}</>
      <title>{app.currentRouteTitle}</title>
    </head>
  );
};
