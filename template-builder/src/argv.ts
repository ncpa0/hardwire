import minimist from "minimist";

type Types = "string" | "number" | "boolean";

type ParsedVars<T extends Record<string, Types>> = {
  [K in keyof T]: T[K] extends "string"
    ? string
    : T[K] extends "number"
    ? number
    : boolean;
};

export class Argv {
  private parsed = minimist(process.argv.slice(2));
  private cmds: Array<{
    name: string;
    required: Record<string, Types>;
    cb: (vars: any, argv: Argv) => void;
  }> = [];

  registerCommand<R extends Record<string, Types>>(
    name: string,
    required: R,
    cb: (vars: ParsedVars<R>, argv: Argv) => void
  ) {
    this.cmds.push({ name, required, cb });
  }

  run() {
    const cmdId = this.parsed._.join(" ");
    const matchingCmd = this.cmds.find((cmd) => cmd.name === cmdId);

    if (!matchingCmd) {
      throw new Error(`Unrecognized command: ${cmdId}`);
    }

    const requiredVars = Object.entries(matchingCmd.required);

    for (const [name, type] of requiredVars) {
      const value = this.parsed[name];
      if (value == null) {
        throw new Error(`Missing required argument: ${name}`);
      }

      if (typeof value !== type) {
        throw new Error(
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
