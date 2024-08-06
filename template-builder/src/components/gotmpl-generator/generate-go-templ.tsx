declare global {
  type AsProxy<T> =
    T extends Array<infer U>
      ? ListProxy<U>
      : T extends object
        ? StructProxy<T>
        : ValueProxy<T>;

  type ListProxy<T> = {
    length: ValueProxy<number>;
    map(
      fn: (item: AsProxy<T>, key: ValueProxy<string>) => JSX.Element,
    ): JSX.Element;
    varname(): string;
  };

  type StructProxy<Struct extends object> = {
    [K in keyof Struct]: AsProxy<Struct[K]>;
  } & {
    varname(): string;
  };

  type ValueProxy<T> = {
    varname(): string;
    toString(): string;
    [Symbol.toHtmlTag](): string;
    [Symbol.toPrimitive](): string;
  };
}

export const structProxy = <T,>(name: string): AsProxy<T> => {
  const toString = () => {
    return `{{${name}}}`;
  };
  const o = {
    name,
    get length() {
      return valueProxy(`len ${name}`);
    },
    varname() {
      return name;
    },
    map<T>(fn: (proxy: AsProxy<T>, key: ValueProxy<string>) => JSX.Element) {
      return <MapArray data={this} render={fn} />;
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

export const valueProxy = <T extends string | number | boolean>(
  value: T,
): ValueProxy<T> => {
  return structProxy(value as any) as any;
};
