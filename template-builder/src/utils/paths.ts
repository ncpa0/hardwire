export const pathCmp = (a: string, b: string): boolean => {
  const aParts = a.split("/").filter(Boolean);
  const bParts = b.split("/").filter(Boolean);

  if (aParts.length !== bParts.length) {
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

export const pathIsWithin = (path1: string, path2: string): boolean => {
  const path1Parts = path1.split("/").filter(Boolean);
  const path2Parts = path2.split("/").filter(Boolean);

  if (path2Parts.length > path1Parts.length) {
    return false;
  }

  for (let i = 0; i < path2Parts.length; i++) {
    if (path1Parts[i] !== path2Parts[i]) {
      return false;
    }
  }

  return true;
};
