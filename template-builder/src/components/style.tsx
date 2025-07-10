import { transform } from "lightningcss";
import esbuild from "esbuild";
import type { ComponentApi } from "jsxte";
import path from "node:path";
import { builderCtx } from "../contexts";
import { escapeHTML } from "bun";

const IS_PROD = process.env.NODE_ENV !== "development";

export type BaseStyleProps = {
  linkAttrs?: JSX.IntrinsicElements["link"];
  styleAttrs?: JSX.IntrinsicElements["style"];
};

export type StyleProps = BaseStyleProps &
  (
    | {
        children?: never;
        package?: never;
        path: string;
        dirname: string;
        inline?: boolean;
        embed?: never;
      }
    | {
        path?: never;
        children?: never;
        package: string;
        inline?: boolean;
        embed?: never;
      }
    | {
        path?: never;
        package?: never;
        children?: JSXTE.TextNodeElement | JSXTE.TextNodeElement[];
        inline?: boolean;
        /**
         * When enabled the contents of this tag will be used.
         */
        embed: true;
      }
  );

const bundleCss = async (filepath: string) => {
  const result = await esbuild.build({
    entryPoints: [filepath],
    bundle: true,
    write: false,
    minify: IS_PROD,
    platform: "browser",
    supported: {
      nesting: false,
    },
    sourcemap: IS_PROD ? undefined : "inline",
  });
  if (result.errors.length > 0) throw new Error(result.errors[0].text);
  return result.outputFiles[0].text;
};

export const Style = async (props: StyleProps, componentApi: ComponentApi) => {
  const extFiles = componentApi.ctx.getOrFail(ExtFilesCtx);
  const builder = componentApi.ctx.getOrFail(builderCtx);

  if (!builder.isBuildPhase) {
    return null;
  }

  const name =
    props.package ?? (props.path ? path.basename(props.path!) : undefined);

  if (name != null && !props.inline) {
    const preBuilt = extFiles.get(name);

    if (preBuilt != null) {
      return (
        <link
          rel="stylesheet"
          href={preBuilt}
          type="text/css"
          media="screen"
          {...(props.linkAttrs ?? {})}
        />
      );
    }
  }

  // const options: MinifyOptions & CompressOptions = {};
  let filepath: string = "";
  let stylesheet: string | undefined;

  if (props.path) {
    filepath = path.join(props.dirname, props.path);
    stylesheet = await bundleCss(filepath);
  } else if (props.package) {
    filepath = props.package!;
    const modulePath = await Bun.resolve(props.package!, builder.entrypointDir);
    // stylesheet = await Bun.file(modulePath).text();
    stylesheet = await bundleCss(modulePath);
  } else if (props.embed) {
    stylesheet = Array.isArray(props.children!)
      ? props.children!.map((n) => n.text).join("\n")
      : props.children!.text;
  }

  if (!stylesheet) {
    return null;
  }

  const result = transform({
    code: Buffer.from(stylesheet) as any,
    filename: filepath,
    minify: IS_PROD,
    sourceMap: !IS_PROD,
  });

  let contents = new TextDecoder().decode(result.code);
  if (filepath) {
    contents = `/* ${filepath} */\n${contents}`;
  }

  if (props.inline) {
    return <style {...(props.styleAttrs ?? {})}>{escapeHTML(contents)}</style>;
  }

  const src = extFiles.register(
    contents,
    props.package ?? (props.path ? path.basename(props.path!) : "inline"),
    "css",
  );

  return (
    <link
      rel="stylesheet"
      href={src}
      type="text/css"
      media="screen"
      {...(props.linkAttrs ?? {})}
    />
  );
};
