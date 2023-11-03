export type StructProxy<Struct extends object> = {
  [K in keyof Struct]: Struct[K] extends object
    ? StructProxy<Struct[K]>
    : ValueProxy<Struct[K]>;
} & {
  varname(): string;
};

export type ValueProxy<T> = {
  varname(): string;
  toString(): string;
};

export const structProxy = <T extends object>(name: string): StructProxy<T> => {
  return new Proxy(
    {
      name,
    },
    {
      get(target, key: string) {
        switch (key) {
          case "varname":
            return () => target.name;
          case "toString":
            return () => `{{${target.name}}}`;
        }
        return structProxy(`${target.name}.${key}`);
      },
    }
  ) as any;
};
