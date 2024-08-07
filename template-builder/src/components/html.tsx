import { HeadProps } from "./head";

export type HtmlProps = JSXTE.PropsWithChildren<
  HeadProps & {
    headContent?: JSX.Element;
  }
>;

export function Html({ children, headContent, ...props }: HtmlProps) {
  return (
    <html>
      <Head {...props}>{headContent}</Head>
      <body hx-ext={props.nomorph ? "" : "morph"}>{children}</body>
    </html>
  );
}
