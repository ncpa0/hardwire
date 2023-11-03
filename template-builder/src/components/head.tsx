export const Head = (props: JSXTE.PropsWithChildren<{}>) => {
  return (
    <head>
      <script
        src="https://unpkg.com/htmx.org@1.9.6"
        integrity="sha384-FhXw7b6AlE/jyjlZH5iHa/tTe9EpJ1Y55RjcgPbjeWMskSxZt1v9qkxLJWNJaGni"
        crossorigin="anonymous"
      ></script>
      {props.children}
    </head>
  );
};
