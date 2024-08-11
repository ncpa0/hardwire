import { HeadProps } from "./head";

export type HtmlProps = JSXTE.PropsWithChildren<{
  headContent?: JSX.Element;
  headProps?: Omit<HeadProps, "extensions">;
  htmlProps?: JSX.IntrinsicElements["html"];
  /**
   * Don't include the idiomorph htmx extension.
   */
  nomorph?: boolean;
  /**
   * Don't include the chunked-transfer htmx extension.
   */
  nochunked?: boolean;
}>;

export function Html({
  children,
  headContent,
  headProps,
  htmlProps,
  nochunked,
  nomorph,
}: HtmlProps) {
  let extensions: string[] = [];
  if (nomorph !== true) {
    extensions.push("morph");
  }
  if (nochunked !== true) {
    extensions.push("chunked-transfer");
  }

  return (
    <html {...htmlProps}>
      <Head {...headProps} extensions={extensions}>
        {headContent}
      </Head>
      <body hx-ext={extensions.join(",")}>{children}</body>
    </html>
  );
}
