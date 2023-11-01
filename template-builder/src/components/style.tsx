import { CompressOptions, MinifyOptions, minify as cssMinify } from "csso";
import type { ComponentApi } from "jsxte";
import path from "node:path";

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

export const Style = async (props: StyleProps, componentApi: ComponentApi) => {
  const options: MinifyOptions & CompressOptions = {};
  let filepath: string;
  let stylesheet: string;

  if (props.path) {
    filepath = path.join(props.dirname, props.path);
    stylesheet = await Bun.file(filepath).text();
  } else {
    filepath = props.package!;
    const modulePath = await Bun.resolve(props.package!, import.meta.dir);
    stylesheet = await Bun.file(modulePath).text();
  }

  const result = await cssMinify(stylesheet, options);

  const contents = `/* ${filepath} */\n${result.css}`;

  if (props.inline) {
    return <style>{contents}</style>;
  }

  const extFiles = componentApi.ctx.getOrFail(ExtFilesCtx);
  const src = extFiles.register(
    contents,
    props.package ?? path.basename(props.path),
    "css"
  );

  return <link rel="stylesheet" href={src} type="text/css" media="screen" />;
};
