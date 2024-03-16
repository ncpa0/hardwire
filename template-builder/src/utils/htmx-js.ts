export const htmxJs = (fn: Function) => {
  const fnstr = String(fn);
  return `javascript: ...(${fnstr})()`;
};
