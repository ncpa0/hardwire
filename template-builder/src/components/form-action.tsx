import { ComponentApi, defineContext } from "jsxte";
import { builderCtx } from "../contexts";
import { defined } from "../utils/defined";
import { IslandMap } from "./island";
import { Client } from "./_helpers.client";
import { arrDedup } from "../utils/arr-dedup";

type BaseActionProps = {
  data?: Record<string, string | number | boolean | ValueProxy<any>>;
  /**
   * List of related island components that should be updated after the action is performed.
   * This field does not override the default action island list but appends to it.
   */
  islands?: JSXTE.Component<any>[];
  /**
   * If the island is an "island list" and this list contains keys, only list items
   * of the specified key(s) will be updated.
   */
  items?: Array<string | ValueProxy<string>>;
  morph?: boolean;
};

export type QuickActionButtonProps = JSX.IntrinsicElements["button"] &
  BaseActionProps & { formProps?: JSX.IntrinsicElements["form"] };
export type FormActionProps = JSX.IntrinsicElements["form"] & BaseActionProps;
export type SubmitActionProps = JSX.IntrinsicElements["button"] & {
  morph?: boolean;
};

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
  items?: Array<string | ValueProxy<string>>;
}>();

let i = 1;
function getNextFormID() {
  return `form_${i++}`;
}

function registerAction(
  api: ComponentApi,
  params: ActionParams,
  islands: string[],
) {
  const builder = api.ctx.getOrFail(builderCtx);
  builder.registerAction({
    resource: params.resource,
    method: params.method,
    action: params.action,
    islandIDs: islands,
  });
}

type ActionParams = {
  resource: string;
  method: HttpMethod;
  action: string;
  /**
   * List of island components that should be updated on the page whenver
   * this action is performed.
   */
  islands?: JSXTE.Component<any>[];
  morph?: boolean;
};

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
export const $action = (actionParams: ActionParams) => {
  const {
    resource,
    method,
    action: actionName,
    islands: relatedIslands = [],
    morph: baseMorph,
  } = actionParams;

  const knownDependantIslands = Array.from(IslandMap.values()).filter(
    (island) => island.resource === resource,
  );

  const baseIslandIDs = arrDedup(
    relatedIslands
      .map((island) => IslandMap.get(island)?.id)
      .concat(knownDependantIslands.map((island) => island.id))
      .filter(defined),
  );

  const uid = getNextFormID();

  const action = {
    get id() {
      return uid;
    },
    QuickButton(
      {
        data,
        islands = [],
        items,
        formProps,
        morph,
        ...props
      }: QuickActionButtonProps,
      api: ComponentApi,
    ) {
      const bldr = api.ctx.getOrFail(builderCtx);

      const btnProps: Record<string, any> = { ...props };
      btnProps["hx-include"] = "#" + uid;
      btnProps["hx-" + method.toLowerCase()] =
        `/__resources/${resource}/actions/${actionName}`;
      btnProps["hx-swap"] = "none";

      let islandsIDs = arrDedup(
        baseIslandIDs.concat(
          islands.map((island) => IslandMap.get(island)?.id).filter(defined),
        ),
      );
      if (islandsIDs.length > 0) {
        const currentPath = bldr.currentRoute.join("/");
        btnProps["hx-headers"] = `javascript: ...${Client.call(
          "formHeaders",
          currentPath,
          islandsIDs,
          (items ?? []).map(String),
          morph ?? baseMorph,
        )}`;
      }

      registerAction(api, actionParams, islandsIDs);

      return (
        <form {...formProps} id={uid}>
          {Object.entries(data ?? {}).map(([key, value]) => {
            return <HiddenInput name={key} value={value} />;
          })}
          <button {...btnProps} />
        </form>
      );
    },
    Form(
      { children, data, islands = [], items, ...props }: FormActionProps,
      api: ComponentApi,
    ) {
      if (api.ctx.has(FormContext)) {
        throw new Error("Form actions cannot be nested.");
      }

      const islandsIDs = arrDedup(
        baseIslandIDs.concat(
          islands.map((island) => IslandMap.get(island)?.id).filter(defined),
        ),
      );

      registerAction(api, actionParams, islandsIDs);

      return (
        <FormContext.Provider
          value={{
            formID: uid,
            islands: islandsIDs,
            items: items,
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
    Submit({ morph, ...props }: SubmitActionProps, api: ComponentApi) {
      const formCtx = api.ctx.getOrFail(FormContext);
      const bldr = api.ctx.getOrFail(builderCtx);

      if (formCtx.formID !== uid) {
        throw new Error(
          "The submit button must be a child of it's own form component.",
        );
      }

      const btnProps: Record<string, any> = { ...props };
      btnProps["hx-include"] = "#" + uid;
      btnProps["hx-" + method.toLowerCase()] =
        `/__resources/${resource}/actions/${actionName}`;
      btnProps["hx-swap"] = "none";

      if (formCtx.islands.length > 0) {
        const currentPath = bldr.currentRoute.join("/");
        btnProps["hx-headers"] = `javascript: ...${Client.call(
          "formHeaders",
          currentPath,
          formCtx.islands,
          (formCtx.items ?? []).map(String),
          morph ?? baseMorph,
        )}`;
      }

      return <button {...btnProps} />;
    },
  };

  return action;
};
