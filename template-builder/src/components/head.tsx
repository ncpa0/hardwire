import { builderCtx } from "../contexts";

export const Head: JSXTE.Component = (props, componentApi) => {
  const app = componentApi.ctx.getOrFail(builderCtx);

  return (
    <head>
      <script
        src="https://unpkg.com/htmx.org@1.9.6"
        integrity="sha384-FhXw7b6AlE/jyjlZH5iHa/tTe9EpJ1Y55RjcgPbjeWMskSxZt1v9qkxLJWNJaGni"
        crossorigin="anonymous"
      ></script>
      <title>{app.currentRouteTitle}</title>
      {props.children}
    </head>
  );
};
