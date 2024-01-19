/**
 * Check if the given path points to a place that is nested within the other given path.
 *
 * @example
 * isSubpath('/foo/bar/baz', '/foo/bar') // true
 * isSubpath('/foo/bar/baz', '/foo/bar/baz') // true
 * isSubpath('/foo/bar', '/foo/bar/baz') // false
 */
export const isSubpath = (nestedPath: string, parentPath: string): boolean => {
  const nestedParts = nestedPath.split("/").filter(Boolean);
  const parentParts = parentPath.split("/").filter(Boolean);

  // a parent will always shorter (e.x. /foo/bar is a parent of /foo/bar/baz, and it's shorter)
  if (nestedParts.length < parentParts.length) {
    return false;
  }

  for (let i = 0; i < nestedParts.length; i++) {
    const aPart = nestedParts[i];
    const bPart = parentParts[i];

    if (bPart == null) {
      break;
    }

    if (aPart.startsWith(":") || bPart.startsWith(":")) {
      continue;
    }

    if (aPart !== bPart) {
      return false;
    }
  }

  return true;
};

/**
 * Check if the two given paths can be the same. Path segments that
 * start with a `:` symbol are treated as wildcards, and don't need
 * to contain the same name for paths to be considered equal.
 *
 * @example
 * '/foo/bar' === '/foo/bar' // true
 * '/foo/:id' === '/foo/bar' // true
 * '/foo/:id/bar' === '/foo/:INDEX/bar' // true
 */
export const pathCompare = (a: string, b: string): boolean => {
  const aParts = a.split("/").filter(Boolean);
  const bParts = b.split("/").filter(Boolean);

  if (aParts.length != bParts.length) {
    return false;
  }

  for (let i = 0; i < aParts.length; i++) {
    const aPart = aParts[i];
    const bPart = bParts[i];

    if (aPart.startsWith(":") && bPart.startsWith(":")) {
      continue;
    }

    if (aPart !== bPart) {
      return false;
    }
  }

  return true;
};
