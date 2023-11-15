import { BuildConfig } from "bun";
import { ComponentApi } from "jsxte";
import fs from "node:fs/promises";
import os from "node:os";
import path from "node:path";
import { builderCtx } from "../contexts";

const IS_PROD = process.env.NODE_ENV !== "development";

const generateRandomName = () => {
  return (
    Math.random().toString(36).slice(2) + Math.random().toString(36).slice(2)
  );
};

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
        package?: never;
        children?: never;
        embed?: never;
        dirname: string;
      }
    | {
        path?: never;
        package: string;
        children?: never;
        embed?: never;
      }
    | {
        path?: never;
        package?: never;
        /**
         * When enabled the contents of this tag will be used.
         */
        embed: true;
        children?: JSXTE.TextNodeElement | JSXTE.TextNodeElement[];
      }
  );

export const Script = async (
  props: ScriptProps,
  componentApi: ComponentApi
) => {
  const extFiles = componentApi.ctx.getOrFail(ExtFilesCtx);
  const builder = componentApi.ctx.getOrFail(builderCtx);

  if (!builder.isBuildPhase) {
    return null;
  }

  const { type = "module", onLoad = () => {}, buildOptions } = props;

  let name: string = "";
  if (props.package) {
    name = props.package;
  } else if (props.path) {
    name = path.basename(props.path!);
  }

  const preBuilt = name ? extFiles.get(name) : undefined;

  if (preBuilt != null) {
    return (
      <script
        type={type === "module" ? "module" : "text/javascript"}
        src={preBuilt}
      />
    );
  }

  const config: BuildConfig = {
    minify: IS_PROD,
    sourcemap: IS_PROD ? "none" : "inline",
    ...buildOptions,
    entrypoints: [],
    target: "browser",
  };

  if (props.path) {
    config.entrypoints = [path.join(props.dirname, props.path)];
  } else if (props.package) {
    const modulePath = await Bun.resolve(props.package!, builder.entrypointDir);
    config.entrypoints = [modulePath];
  } else if (props.embed && props.children) {
    const tmpdir = await fs.mkdtemp(
      path.join(os.tmpdir(), "template-builder-js")
    );
    const tmpFile = path.join(tmpdir, `${generateRandomName()}.ts`);
    await Bun.write(
      tmpFile,
      Array.isArray(props.children)
        ? props.children.map((n) => n.text).join("\n")
        : props.children.text
    );
    console.debug(props.children);
    config.entrypoints = [tmpFile];
  }

  if (config.entrypoints.length === 0) {
    return <></>;
  }

  const result = await Bun.build(config);

  if (!result.success) {
    throw new Error(`Build failed. [${config.entrypoints[0]}]`);
  }

  let contents = await result.outputs[0].text().then((t) => t.trim());
  if (props.path) {
    contents = `/* ${props.path} */\n${contents}`;
  } else if (props.package) {
    contents = `/* ${props.package} */\n${contents}`;
  }

  const tmp = onLoad(contents);
  if (tmp) {
    contents = tmp;
  }

  const src = extFiles.register(contents, name, "js");

  return (
    <script type={type === "module" ? "module" : "text/javascript"} src={src} />
  );
};
