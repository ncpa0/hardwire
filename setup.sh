#!/bin/sh

mkdir pages
cd pages

bun add -D jsxte

touch index.d.ts
touch tsconfig.json

dts="import \"jsxte\";

declare global {
  type ScriptPropsBase = {
    /**
     * When module, it will be bundled into an ESM format and imported using the script tag with a
     * module type.
     *
     * When iife, it will be bundled into an IIFE format and imported using the script tag with a
     * text/javascript type.
     *
     * When global, it will be bundled into a CJS format, ran as IIFE and all exports will be assigned
     * to the global object.
     */
    type?: \"module\";
    onLoad?: (contents: string) => string | undefined;
    buildOptions?: BuildConfig;
  };

  type ScriptProps = ScriptPropsBase &
    (
      | {
          path: string;
          dirname: string;
          package?: never;
        }
      | {
          path?: never;
          package: string;
        }
    );

  const Script: (
    props: ScriptProps,
    componentApi: ComponentApi
  ) => Promise<JSX.SyncElement>;

  export type StyleProps =
    | {
        path: string;
        dirname: string;
        package?: never;
        inline?: boolean;
      }
    | {
        path?: never;
        package: string;
        inline?: boolean;
      };

  const Style: (
    props: StyleProps,
    componentApi: ComponentApi
  ) => Promise<JSX.SyncElement>;

  export function Router(
    props: JSXTE.PropsWithChildren<{ id: string }>
  ): JSX.Element;
  export function Route(
    props: JSXTE.PropsWithChildren<{
      path: string;
    }>
  ): JSX.Element;
  export function Link(
    props: JSXTE.PropsWithChildren<
      JSX.IntrinsicElements[\"a\"] & {
        href: string;
      }
    >
  ): JSX.Element;
}"

tsconf="{
    \"compilerOptions\": {
        \"jsx\": \"react-jsx\",
        \"jsxImportSource\": \"jsxte\"
    }
}"

echo "$dts" > htmx-framework.d.ts
echo "$tsconf" > tsconfig.json