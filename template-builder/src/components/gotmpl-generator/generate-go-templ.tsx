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
  [Symbol.toHtmlTag](): string;
  [Symbol.toPrimitive](): string;
};

export const structProxy = <T extends object>(name: string): StructProxy<T> => {
  const toString = () => {
    return `{{${name}}}`;
  };
  const o = {
    name,
    varname: () => {
      return name;
    },
    toString: toString,
    [Symbol.toHtmlTag]: toString,
    [Symbol.toPrimitive]: toString,
  };
  const ownKeys: Array<string | symbol> = [
    ...Object.getOwnPropertyNames(o),
    ...Object.getOwnPropertySymbols(o),
  ];
  return new Proxy(o, {
    get(target, key: string | symbol) {
      if (key in target) {
        // @ts-expect-error
        return target[key];
      }
      return structProxy(`${target.name}.${key as string}`);
    },
    ownKeys() {
      return ownKeys;
    },
    has(_, p) {
      return ownKeys.includes(p as string);
    },
  }) as any;
};
