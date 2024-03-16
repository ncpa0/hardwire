import { DynamicFragmentProps } from "./dynamic-fragment";

export const IslandMap = new Map<
  JSXTE.Component<any>,
  {
    id: string;
    fragmentID: string;
  }
>();

export function $island<T extends any, P extends object = {}>(
  options: Omit<DynamicFragmentProps<T>, "render"> & { id: string },
  Component: (props: P, data: AsProxy<T>) => JSX.Element
): JSXTE.Component<P> {
  const { id, ...dynamicFragmentProps } = options;

  const islandEntry = {
    id,
    fragmentID: "",
  };

  const island: JSXTE.Component<P> = (props) => {
    return (
      <div id={id}>
        <DynamicFragment
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
      </div>
    );
  };

  IslandMap.set(island, islandEntry);

  return island;
}
