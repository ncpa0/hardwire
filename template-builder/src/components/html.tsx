import { builderCtx } from "../contexts";
import { SetupClientHelpers } from "./_helpers.client";
import { HeadProps as InternalHeadProps } from "./head";
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
    extensions.push("htmx-ext-head-support");
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
   * The version of htmx to use.
   *
   * @default "2.0.1"
   */
  htmxVersion?: string;
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
  const htmxVer = props.htmxVersion ?? "2.0.1";
  let integrity: string | undefined = undefined;

  switch (htmxVer) {
    case "2.0.1":
      integrity =
        "sha384-QWGpdj554B4ETpJJC9z+ZHJcA/i59TyjxEPXiiUgN2WmTyV5OEZWCD6gQhgkdpB/";
      break;
    case "1.9.12":
      integrity =
        "sha384-ujb1lZYygJmzgSwoxRggbCHcjc0rB2XoQrxeTUQyRjrOnlCoYta87iKBWq3EsdM2";
      break;
    case "1.9.11":
      integrity =
        "sha384-0gxUXCCR8yv9FM2b+U3FDbsKthCI66oH5IA9fHppQq9DDMHuMauqq1ZHBpJxQ0J0";
      break;
    case "1.9.10":
      integrity =
        "sha384-D1Kt99CQMDuVetoL1lrYwg5t+9QdHe7NLX/SoJYkXDFfX37iInKRy5xLSi8nO7UC";
      break;
    case "1.9.9":
      integrity =
        "sha384-QFjmbokDn2DjBjq+fM+8LUIVrAgqcNW2s0PjAxHETgRn9l4fvX31ZxDxvwQnyMOX";
      break;
    case "1.9.8":
      integrity =
        "sha384-rgjA7mptc2ETQqXoYC3/zJvkU7K/aP44Y+z7xQuJiVnB/422P/Ak+F/AqFR7E4Wr";
      break;
    case "1.9.7":
      integrity =
        "sha384-EAzY246d6BpbWR7sQ8+WEm40J8c3dHFsqC58IgPlh4kMbRRI6P6WA+LA/qGAyAu8";
      break;
    case "1.9.6":
      integrity =
        "sha384-FhXw7b6AlE/jyjlZH5iHa/tTe9EpJ1Y55RjcgPbjeWMskSxZt1v9qkxLJWNJaGni";
      break;
    case "1.9.5":
      integrity =
        "sha384-xcuj3WpfgjlKF+FXhSQFQ0ZNr39ln+hwjN3npfM9VBnUskLolQAcN80McRIVOPuO";
      break;
    case "1.9.4":
      integrity =
        "sha384-zUfuhFKKZCbHTY6aRR46gxiqszMk5tcHjsVFxnUo8VMus4kHGVdIYVbOYYNlKmHV";
      break;
    case "1.9.3":
      integrity =
        "sha384-lVb3Rd/Ca0AxaoZg5sACe8FJKF0tnUgR2Kd7ehUOG5GCcROv5uBIZsOqovBAcWua";
      break;
    case "1.9.2":
      integrity =
        "sha384-L6OqL9pRWyyFU3+/bjdSri+iIphTN/bvYyM37tICVyOJkWZLpP2vGn6VUEXgzg6h";
      break;
    case "1.9.1":
      integrity =
        "sha384-KReoNuwj58fe4zgWyjj5a1HrvXYPBeV0a3bNPVjK7n5FdsGC41fHRx6sq5tONeP0";
      break;
    case "1.9.0":
      integrity =
        "sha384-aOxz9UdWG0yBiyrTwPeMibmaoq07/d3a96GCbb9x60f3mOt5zwkjdbcHFnKH8qls";
      break;
  }

  if (props.htmxIntegrityHash) {
    integrity = props.htmxIntegrityHash;
  }

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
      <script
        src={`https://unpkg.com/htmx.org@${htmxVer}`}
        integrity={integrity}
        crossorigin="anonymous"
      />
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
