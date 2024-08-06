#!/usr/bin/env bun

/// <reference types="bun-types" />
import path from "node:path";
import { Argv } from "./argv";

declare global {
  namespace JSXTE {
    interface AttributeAcceptedTypes {
      ALL: ValueProxy<any>;
    }
  }
}

async function main() {
  const argv = new Argv("hardwire-html-generator");

  const init = argv.registerCommand(
    "init",
    {
      dir: "string",
    },
    async (args) => {
      const { initCmd } = await import("./command-init/init");
      const wd = path.resolve(args.dir);
      await initCmd(wd);
    },
  );

  const build = argv.registerCommand(
    "build",
    {
      src: "string",
      outdir: "string",
      staticdir: "string",
      staticurl: "string",
    },
    async (args) => {
      const { buildCmd } = await import("./command-build/build");
      const srcFile = path.resolve(args.src);
      const outDir = path.resolve(args.outdir);
      const staticDir = path.resolve(args.staticdir);
      await buildCmd(srcFile, outDir, staticDir, args.staticurl);
    },
  );

  init.setDescription("Initializes a new project in the given directory.");
  build.setDescription(
    "Generates static HTML files from the given JSX component.",
  );

  argv.run();
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
