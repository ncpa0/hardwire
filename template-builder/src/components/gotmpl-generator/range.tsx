import { ComponentApi } from "jsxte";
import { StructProxy, ValueProxy, structProxy } from "./generate-go-templ";

type ProxyFor<T> = T extends object ? StructProxy<T> : ValueProxy<T>;

type MapProps<Data extends StructProxy<any[]>> = {
  data: Data;
  render: (data: Data[number]) => JSX.Element;
};

export const MapArray = async <Data extends StructProxy<any[]>>(
  props: MapProps<Data>,
  compApi: ComponentApi
) => {
  return `{{range ${props.data.varname()}}} ${await compApi.renderAsync(
    props.render(structProxy(""))
  )} {{end}}`;
};
