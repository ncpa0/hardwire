import { ComponentApi } from "jsxte";
import { render } from "../../renderer";

type Primitive = number | string | boolean;

const TEMPL_QUOTE = "@#34T;";

function quote(str: string): string {
  return `${TEMPL_QUOTE}${str.replaceAll('"', TEMPL_QUOTE)}${TEMPL_QUOTE}`;
}

function templateValue(value: Primitive | ValueProxy<any>): string {
  if (typeof value === "object" && "varname" in value) {
    return value.varname();
  }
  switch (typeof value) {
    case "string":
      return quote(value);
    case "number":
      return String(value);
    case "boolean":
      return value ? "true" : "false";
  }
}

type IfProps = JSXTE.PropsWithChildren<{
  condition: (conditionBuilder: ConditionBuilder) => ConditionBuilder;
  negate?: boolean;
  then: () => JSX.Element;
  else?: () => JSX.Element;
}>;

export const If = async (props: IfProps, comApi: ComponentApi) => {
  const forTrue = await render(<>{props.then()}</>, comApi);
  const forFalse = props.else
    ? " {{else}}" + (await render(<>{props.else()}</>, comApi))
    : "";

  return `{{if${props.negate ? " not " : ""} ${props
    .condition(new ConditionBuilder())
    .varname()}}}\n${forTrue}${forFalse}\n{{end}}`;
};

type Comparable = string | number;
type ComplexCondition = ConditionBuilder | ValueProxy<boolean>;

export class ConditionBuilder {
  protected content: string = "";

  varname(): string {
    return this.content;
  }

  not(condition: ComplexCondition): ConditionBuilder {
    const next = new ConditionBuilder();
    next.content = `not ${condition.varname()}`;
    return next;
  }
  and(...conditions: ComplexCondition[]): ConditionBuilder {
    const next = new ConditionBuilder();
    next.content = `${conditions.at(-1)!.varname()}`;
    for (let i = conditions.length - 2; i >= 0; i--) {
      const c = conditions[i]!;
      next.content = `(and ${c.varname()} ${next.content})`;
    }
    return next;
  }
  or(...conditions: ComplexCondition[]): ConditionBuilder {
    const next = new ConditionBuilder();
    next.content = `${conditions.at(-1)!.varname()}`;
    for (let i = conditions.length - 2; i >= 0; i--) {
      const c = conditions[i]!;
      next.content = `(or ${c.varname()} ${next.content})`;
    }
    return next;
  }
  /** Equals to */
  equal(
    a: ValueProxy<any> | Primitive,
    b: ValueProxy<any> | Primitive
  ): ConditionBuilder {
    const next = new ConditionBuilder();
    next.content = `(eq ${templateValue(a)} ${templateValue(b)})`;
    return next;
  }
  /** Not equals to */
  notEqual(
    a: ValueProxy<any> | Primitive,
    b: ValueProxy<any> | Primitive
  ): ConditionBuilder {
    const next = new ConditionBuilder();
    next.content = `(ne ${templateValue(a)} ${templateValue(b)})`;
    return next;
  }
  /** Less than */
  lessThan<T extends Comparable>(
    a: ValueProxy<T> | T,
    b: ValueProxy<T> | T
  ): ConditionBuilder {
    const next = new ConditionBuilder();
    next.content = `(lt ${templateValue(a)} ${templateValue(b)})`;
    return next;
  }
  /** Less than or equals to */
  lessEqualThan<T extends Comparable>(
    a: ValueProxy<T> | T,
    b: ValueProxy<T> | T
  ): ConditionBuilder {
    const next = new ConditionBuilder();
    next.content = `(le ${templateValue(a)} ${templateValue(b)})`;
    return next;
  }
  /** Greater than */
  greaterThan<T extends Comparable>(
    a: ValueProxy<T> | T,
    b: ValueProxy<T> | T
  ): ConditionBuilder {
    const next = new ConditionBuilder();
    next.content = `(gt ${templateValue(a)} ${templateValue(b)})`;
    return next;
  }
  /** Greater than or equals to */
  greaterEqualThan<T extends Comparable>(
    a: ValueProxy<T> | T,
    b: ValueProxy<T> | T
  ): ConditionBuilder {
    const next = new ConditionBuilder();
    next.content = `(ge ${templateValue(a)} ${templateValue(b)})`;
    return next;
  }
  /**
   * Check if the type of the value is the given type. The type is a Go type name.
   */
  typeofIs(value: ValueProxy<any>, type: string): ConditionBuilder {
    const next = new ConditionBuilder();
    next.content = `(eq ${quote(type)} (printf ${quote("%T")} ${value.varname()}))`;
    return next;
  }
}
