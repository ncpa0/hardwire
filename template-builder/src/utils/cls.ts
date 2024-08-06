export function cls(
  ...args: Array<string | ValueProxy<string> | undefined | null>
) {
  let cname = "";
  for (let i = 0; i < args.length; i++) {
    const name = args[i];
    if (name) {
      cname += " " + name.toString().trim();
    }
  }
  return cname;
}
