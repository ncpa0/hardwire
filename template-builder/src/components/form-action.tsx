import { ComponentApi, defineContext } from "jsxte";
import { builderCtx } from "../contexts";
import { defined } from "../utils/defined";
import { htmxJs } from "../utils/htmx-js";
import { IslandMap } from "./island";

export type FormActionProps = JSX.IntrinsicElements["form"] & {
  data?: Record<string, string | number | boolean | ValueProxy<any>>;
  /**
   * List of related island components that should be updated after the action is performed.
   * This field does not override the default action island list but appends to it.
   */
  islands?: JSXTE.Component<any>[];
};
export type SubmitActionProps = JSX.IntrinsicElements["button"];

export type HttpMethod = "GET" | "POST" | "PUT" | "DELETE" | "PATCH";

function HiddenCheckbox({ name, value }: { name: string; value: boolean }) {
  if (value) {
    return (
      <input
        style="display: none;"
        type="checkbox"
        name={name}
        checked="true"
      />
    );
  } else {
    return <input style="display: none;" type="checkbox" name={name} />;
  }
}

function HiddenInput({
  name,
  value,
}: {
  name: string;
  value: string | number | boolean | ValueProxy<any>;
}) {
  if (typeof value === "object" && "varname" in value) {
    return (
      <If
        condition={(c) => c.typeofIs(value, "boolean")}
        then={() => {
          return (
            <If
              condition={(c) => c.equal(value, true)}
              then={() => <HiddenCheckbox name={name} value={true} />}
              else={() => <HiddenCheckbox name={name} value={false} />}
            />
          );
        }}
        else={() => <input type="hidden" name={name} value={value} />}
      />
    );
  }
  if (typeof value === "boolean") {
    return <HiddenCheckbox name={name} value={value} />;
  }
  return <input type="hidden" name={name} value={String(value)} />;
}

const FormContext = defineContext<{
  formID: string;
  islands: string[];
}>();

let i = 1;
function getNextFormID() {
  return `form_${i++}`;
}

/**
 * @param method - The HTTP method to use for the form submission
 * @param action - The action name to trigger on the server
 * @param relatedIslands - List of related island components that should be updated after the action is performed
 *
 * @example
 * const action = createFormAction("POST", "create-article");
 *
 * <action.Form>
 *  <input name="Title" />
 *  <input name="Body" />
 *  <action.Submit>Submit</action.Submit>
 * </action.Form>
 *
 * // clicking the submit, will send a POST request that will trigger
 * // the "create-article" action that's registered on the server.
 */
export const createFormAction = (
  method: HttpMethod,
  action: string,
  relatedIslands: JSXTE.Component<any>[] = []
) => {
  const uid = getNextFormID();

  return {
    get id() {
      return uid;
    },
    Form: (
      { children, data, islands = [], ...props }: FormActionProps,
      api: ComponentApi
    ) => {
      if (api.ctx.has(FormContext)) {
        throw new Error("Form actions cannot be nested.");
      }

      const islandsIDs = relatedIslands
        .concat(islands)
        .map((island) => IslandMap.get(island)?.id)
        .filter(defined);

      return (
        <FormContext.Provider
          value={{
            formID: uid,
            islands: islandsIDs,
          }}
        >
          <form {...props} id={uid}>
            {Object.entries(data ?? {}).map(([key, value]) => {
              return <HiddenInput name={key} value={value} />;
            })}
            {children}
          </form>
        </FormContext.Provider>
      );
    },
    Submit: (props: SubmitActionProps, api: ComponentApi) => {
      const formCtx = api.ctx.getOrFail(FormContext);
      const bldr = api.ctx.getOrFail(builderCtx);

      if (formCtx.formID !== uid) {
        throw new Error(
          "The submit button must be a child of it's own form component."
        );
      }

      const btnProps: Record<string, any> = { ...props };
      btnProps["hx-include"] = "#" + uid;
      btnProps["hx-" + method.toLowerCase()] = `/__actions/${action}`;
      btnProps["hx-swap"] = "none";

      if (formCtx.islands.length > 0) {
        btnProps["hx-headers"] = htmxJs(function () {
          const relatedIslands = ["%ISLANDS%"];
          const presentIslands = [];
          for (let i = 0; i < relatedIslands.length; i++) {
            const islandID = relatedIslands[i];
            const islandEl = document.getElementById(islandID);
            if (islandEl) {
              presentIslands.push(islandID);
            }
          }
          return {
            "Hardwire-Islands-Update": presentIslands.join(";"),
            "Hardwire-Dynamic-Fragment-Request": "%DYN_PATHNAME%",
          };
        })
          .replace('["%ISLANDS%"]', JSON.stringify(formCtx.islands))
          .replace("%DYN_PATHNAME%", "/" + bldr.currentRoute.join("/"));
      }

      return <button {...btnProps} />;
    },
  };
};
