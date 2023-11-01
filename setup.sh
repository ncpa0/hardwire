#!/bin/sh

mkdir pages
cd pages

bun add -D jsxte

touch index.d.ts
touch tsconfig.json

dts="import \"jsxte\";

declare global {
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
}

export {};"

tsconf="{
    \"compilerOptions\": {
        \"jsx\": \"react-jsx\",
        \"jsxImportSource\": \"jsxte\"
    }
}"

echo "$dts" > index.d.ts
echo "$tsconf" > tsconfig.json