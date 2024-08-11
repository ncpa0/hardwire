import { DynamicFragmentProps } from "./dynamic-fragment";

export type IslandDefinition = {
  id: string;
  fragmentID: string;
  resource: string;
  type: "list" | "basic";
};

export const IslandMap = new Map<JSXTE.Component<any>, IslandDefinition>();

export function $island<T extends any, P extends object = {}>(
  options: Omit<DynamicFragmentProps<T>, "render"> & { id: string },
  Component: (props: P, data: AsProxy<T>) => JSX.Element,
): JSXTE.Component<P> {
  const { id, ...dynamicFragmentProps } = options;

  const islandEntry: IslandDefinition = {
    id,
    fragmentID: "",
    resource: dynamicFragmentProps.require,
    type: "basic",
  };

  const island: JSXTE.Component<P> = (props) => {
    return (
      <DynamicFragment
        id={id}
        {...dynamicFragmentProps}
        render={(data: AsProxy<any>) => {
          // @ts-expect-error
          return Component(props, data);
        }}
        // @ts-expect-error
        __fragidgetter={(id: string) => {
          islandEntry.fragmentID = id;
        }}
      />
    );
  };

  IslandMap.set(island, islandEntry);

  return island;
}
