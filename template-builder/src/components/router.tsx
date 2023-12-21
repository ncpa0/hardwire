import { ComponentApi } from "jsxte";
import { builderCtx, routerCtx } from "../contexts";

export const Switch = (
  props: JSXTE.PropsWithChildren<{ id: string } & JSX.IntrinsicElements["div"]>,
  compApi: ComponentApi
) => {
  const { children, ...rest } = props;
  const app = compApi.ctx.getOrFail(builderCtx);
  app.addRouter(props.id);
  return (
    <div {...rest}>
      <routerCtx.Provider
        value={{
          containerID: rest.id,
        }}
      >
        {children}
      </routerCtx.Provider>
    </div>
  );
};
