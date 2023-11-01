import { BuildConfig } from "bun";
import { ComponentApi } from "jsxte";
import path from "node:path";
import { builderCtx } from "../contexts";

const IS_PROD = process.env.NODE_ENV !== "development";

type ScriptPropsBase = {
  /**
   * When module, it will be bundled into an ESM format and imported using the script tag with a
   * `module` type.
   *
   * When iife, it will be bundled into an IIFE format and imported using the script tag with a
   * `text/javascript` type.
   *
   * When global, it will be bundled into a CJS format, ran as IIFE and all exports will be assigned
   * to the global object.
   */
  type?: "module";
  onLoad?: (contents: string) => string | undefined;
  buildOptions?: BuildConfig;
};

export type ScriptProps = ScriptPropsBase &
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

export const Script = async (
  props: ScriptProps,
  componentApi: ComponentApi
) => {
  const builder = componentApi.ctx.getOrFail(builderCtx);

  if (!builder.isBuildStep) {
    return null;
  }

  const { type = "module", onLoad = () => {}, buildOptions } = props;

  const config: BuildConfig = {
    minify: IS_PROD,
    sourcemap: IS_PROD ? "none" : "inline",
    ...buildOptions,
    entrypoints: [],
    target: "browser",
  };

  if (props.path) {
    config.entrypoints = [path.join(props.dirname, props.path)];
  } else {
    const modulePath = await Bun.resolve(props.package!, builder.entrypointDir);
    config.entrypoints = [modulePath];
  }

  const result = await Bun.build(config);

  if (!result.success) {
    throw new Error(`Build failed. [${config.entrypoints[0]}]`);
  }

  let contents = await result.outputs[0].text().then((t) => t.trim());

  const tmp = onLoad(contents);
  if (tmp) {
    contents = tmp;
  }

  const extFiles = componentApi.ctx.getOrFail(ExtFilesCtx);
  const src = extFiles.register(
    contents,
    props.package ?? path.basename(props.path),
    "js"
  );

  return (
    <script type={type === "module" ? "module" : "text/javascript"} src={src} />
  );
};
