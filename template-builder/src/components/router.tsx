import { ComponentApi } from "jsxte";
import { builderCtx, routerCtx } from "../contexts";

export const Switch = (
  props: JSXTE.PropsWithChildren<{ id: string }>,
  compApi: ComponentApi
) => {
  const app = compApi.ctx.getOrFail(builderCtx);
  app.addRouter(props.id);
  return (
    <div id={props.id}>
      <routerCtx.Provider
        value={{
          containerID: props.id,
        }}
      >
        {props.children}
      </routerCtx.Provider>
    </div>
  );
};
