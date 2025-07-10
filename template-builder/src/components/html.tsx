import { builderCtx } from "../contexts";
import { SetupClientHelpers } from "./_helpers.client";
import { ComponentApi, defineContext } from "jsxte";

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
      <InternalHead {...headProps} extensions={extensions}>
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

  return (
    <head>
      <meta charset="utf-8" />
      <meta http-equiv="x-ua-compatible" content="IE=edge" />
      <meta name="viewport" content="width=device-width, initial-scale=1" />
      <Script package="htmx.org" />
      {props.extensions?.map((ext) => {
        switch (ext) {
          case "chunked-transfer": {
            return <Script package="htmx.ext...chunked-transfer" />;
            break;
          }
          case "morph": {
            return <Script package="idiomorph" />;
            break;
          }
          case "head-support": {
            return <Script package="htmx-ext-head-support" />;
            break;
          }
        }
      })}
      <script src={hwScript} />
      <>{props.children}</>
      <title>{app.currentRouteTitle}</title>
    </head>
  );
};
