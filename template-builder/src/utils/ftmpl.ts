export function ftmpl(strs: ReadonlyArray<string>, ...params: any[]): string {
  let result = "";
  for (let i = 0; i < strs.length; i++) {
    result += strs[i];
    if (i < params.length) {
      result += params[i];
    }
  }
  return result.trim() + "\n";
}
