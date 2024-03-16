import { ComponentApi, defineContext } from "jsxte";
import { render } from "../../renderer";
import { structProxy, valueProxy } from "./generate-go-templ";

type MapProps<T> = {
  data: ListProxy<T>;
  render: (data: AsProxy<T>, key: ValueProxy<string>) => JSX.Element;
};

const MapArrayArgNameContext = defineContext<{
  i: number;
}>();

const getNextName = (i: number) => {
  // return `$_${Math.random().toString(36).slice(2)}`;
  return `$range_elem${i}`;
};

export const MapArray = async <T,>(
  props: MapProps<T>,
  compApi: ComponentApi
) => {
  let { i = 1 } = compApi.ctx.get(MapArrayArgNameContext) ?? {};

  const name = getNextName(i);
  const keyName = `${name}_key`;

  compApi.ctx.set(MapArrayArgNameContext, { i: i + 1 });

  return `{{range ${keyName}, ${name} := ${props.data.varname()}}}\n${await render(
    props.render(structProxy(name), valueProxy(keyName)),
    compApi
  )}\n{{end}}`;
};
