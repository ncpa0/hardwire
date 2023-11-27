import { ComponentApi } from "jsxte";
import {
  StructProxy,
  ValueProxy,
  structProxy,
  valueProxy,
} from "./generate-go-templ";

type MapProps<Data extends StructProxy<any[]>> = {
  data: Data;
  render: (data: Data[number], key: ValueProxy<string>) => JSX.Element;
};

const randName = () => {
  return `$_${Math.random().toString(36).slice(2)}`;
};

export const MapArray = async <Data extends StructProxy<any[]>>(
  props: MapProps<Data>,
  compApi: ComponentApi
) => {
  const name = randName();
  const keyName = `${name}_key`;

  return `{{range ${keyName}, ${name} := ${props.data.varname()}}}\n${await compApi.renderAsync(
    props.render(structProxy(name), valueProxy(keyName))
  )}\n{{end}}`;
};
