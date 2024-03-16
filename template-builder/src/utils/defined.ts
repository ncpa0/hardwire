export function defined<T>(value: T | undefined | null): value is T {
  return value != null;
}
