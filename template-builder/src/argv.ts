import minimist from "minimist";

type Types = "string" | "number" | "boolean";

type ParsedVars<T extends Record<string, Types>> = {
  [K in keyof T]: T[K] extends "string"
    ? string
    : T[K] extends "number"
    ? number
    : boolean;
};

type CmdEntry = {
  name: string;
  required: Record<string, Types>;
  cb: (vars: any, argv: Argv) => void;
  description?: string;
};

export class Argv {
  private readonly parsed = minimist(process.argv.slice(2));
  private readonly cmds: Array<CmdEntry> = [];

  constructor(private readonly scriptName: string) {
    Object.freeze(this.parsed);
    Object.freeze(this.parsed._);
  }

  private error(msg: string) {
    console.error(msg);
    process.exit(1);
  }

  private printHelp() {
    const out = Bun.stdout.writer();

    if (this.cmds.length) {
      out.write(`USAGE: ${this.scriptName} <command> [args]\n\n`);
    } else {
      out.write(`USAGE: ${this.scriptName} [args]\n\n`);
    }

    if (this.cmds.length) {
      const longestCmdName = Math.max(
        ...this.cmds.map((cmd) => cmd.name.length)
      );

      out.write("COMMANDS:\n");
      for (const cmd of this.cmds) {
        out.write(`  ${cmd.name}`);

        if (cmd.description) {
          const spacing = " ".repeat(longestCmdName - cmd.name.length + 6);
          out.write(`${spacing}${cmd.description}`);
        }

        out.write("\n");
      }
      out.write("\n");
    }

    out.flush();
  }

  registerCommand<R extends Record<string, Types>>(
    name: string,
    required: R,
    cb: (vars: ParsedVars<R>, argv: Argv) => void
  ) {
    const entry: CmdEntry = { name, required, cb };
    this.cmds.push(entry);

    return {
      setDescription(description: string) {
        entry.description = description;
        return this;
      },
    };
  }

  run() {
    if (this.parsed.help) {
      return this.printHelp();
    }

    const cmdId = this.parsed._.join(" ");
    const matchingCmd = this.cmds.find((cmd) => cmd.name === cmdId);

    if (!matchingCmd) {
      return this.printHelp();
    }

    const requiredVars = Object.entries(matchingCmd.required);

    for (const [name, type] of requiredVars) {
      const value = this.parsed[name];
      if (value == null) {
        return this.error(`--${name} was not specified`);
      }

      if (typeof value !== type) {
        return this.error(
          `Invalid type for argument ${name}: expected ${type}, got ${typeof value}`
        );
      }
    }

    const parsedVars = Object.fromEntries(
      requiredVars.map(([name, type]) => [name, this.parsed[name]])
    );

    return matchingCmd.cb(parsedVars, this);
  }
}
