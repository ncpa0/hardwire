import { randomUUID } from "crypto";

export type FormActionProps = JSX.IntrinsicElements["form"];
export type SubmitActionProps = JSX.IntrinsicElements["button"];

export type HttpMethod = "GET" | "POST" | "PUT" | "DELETE" | "PATCH";

/**
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
export const createFormAction = (method: HttpMethod, action: string) => {
  const uid = "form_" + randomUUID().replaceAll("-", "").substring(0, 8);

  return {
    get id() {
      return uid;
    },
    Form: (props: FormActionProps) => {
      return <form {...props} id={uid}></form>;
    },
    Submit: (props: SubmitActionProps) => {
      const btnProps: Record<string, any> = { ...props };
      btnProps["hx-include"] = "#" + uid;
      btnProps["hx-" + method.toLowerCase()] = `/__actions/${action}`;

      return <button {...btnProps} />;
    },
  };
};
