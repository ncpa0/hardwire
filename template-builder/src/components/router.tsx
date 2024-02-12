import { ComponentApi } from "jsxte";
import { builderCtx, routerCtx } from "../contexts";

const SW_CLNAME = "__route-switch";

export const Switch = (
  props: JSXTE.PropsWithChildren<{ id: string } & JSX.IntrinsicElements["div"]>,
  compApi: ComponentApi
) => {
  const { children, class: cln, ...rest } = props;
  const app = compApi.ctx.getOrFail(builderCtx);
  app.addRouter(props.id);
  return (
    <div {...rest} class={cln ? cln + ` ${SW_CLNAME}` : SW_CLNAME}>
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
