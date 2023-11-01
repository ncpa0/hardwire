import { CompressOptions, MinifyOptions, minify as cssMinify } from "csso";
import type { ComponentApi } from "jsxte";
import path from "node:path";
import { builderCtx } from "../contexts";

export type StyleProps =
  | {
      children?: never;
      package?: never;
      path: string;
      dirname: string;
      inline?: boolean;
    }
  | {
      path?: never;
      children?: never;
      package: string;
      inline?: boolean;
    }
  | {
      path?: never;
      package?: never;
      children: string | string[];
      inline?: boolean;
    };

export const Style = async (props: StyleProps, componentApi: ComponentApi) => {
  const extFiles = componentApi.ctx.getOrFail(ExtFilesCtx);
  const builder = componentApi.ctx.getOrFail(builderCtx);

  if (!builder.isBuildPhase) {
    return null;
  }

  const name =
    props.package ?? props.path ? path.basename(props.path!) : "inline";

  if (name !== "inline" && !props.inline) {
    const preBuilt = extFiles.get(name);

    if (preBuilt != null) {
      return (
        <link rel="stylesheet" href={preBuilt} type="text/css" media="screen" />
      );
    }
  }

  const options: MinifyOptions & CompressOptions = {};
  let filepath: string;
  let stylesheet: string;

  if (props.path) {
    filepath = path.join(props.dirname, props.path);
    stylesheet = await Bun.file(filepath).text();
  } else if (props.package) {
    filepath = props.package!;
    const modulePath = await Bun.resolve(props.package!, builder.entrypointDir);
    stylesheet = await Bun.file(modulePath).text();
  } else {
    filepath = "inline";
    stylesheet = Array.isArray(props.children!)
      ? props.children!.join("\n")
      : props.children!;
  }

  const result = await cssMinify(stylesheet, options);

  const contents = `/* ${filepath} */\n${result.css}`;

  if (props.inline) {
    return <style>{contents}</style>;
  }

  const src = extFiles.register(
    contents,
    props.package ?? props.path ? path.basename(props.path!) : "inline",
    "css"
  );

  return <link rel="stylesheet" href={src} type="text/css" media="screen" />;
};
