export const sanitizeForHtml = (value: string): string => {
  return value
    .replaceAll('"', "&quot;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;");
};
