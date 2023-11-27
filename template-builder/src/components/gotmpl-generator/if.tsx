import { ComponentApi, defineContext } from "jsxte";
import { ValueProxy } from "./generate-go-templ";

type IfProps = JSXTE.PropsWithChildren<{
  condition: (conditionBuilder: ConditionBuilder) => ConditionBuilder;
  negate?: boolean;
}>;

const ifCtx = defineContext<{ onElse(templ: string): void }>();

export const If = async (props: IfProps, comApi: ComponentApi) => {
  let elseTempl: string = "";
  const onElse = (templ: string) => {
    elseTempl = ` {{else}}\n${templ}\n`;
  };

  const forTrue = await comApi.renderAsync(
    <ifCtx.Provider value={{ onElse }}>{props.children}</ifCtx.Provider>
  );

  return `{{if${props.negate ? " not " : ""} ${props
    .condition(new ConditionBuilder())
    .varname()}}}\n${forTrue}${elseTempl}\n{{end}}`;
};

export const Else = async (
  props: JSXTE.PropsWithChildren<{}>,
  comApi: ComponentApi
) => {
  const ctx = comApi.ctx.getOrFail(ifCtx);
  ctx.onElse(await comApi.renderAsync(<>{props.children}</>));
  return null;
};

type Comparable = string | number;
type ComplexCondition = ConditionBuilder | ValueProxy<boolean>;

class ConditionBuilder {
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
    // @ts-expect-error
    next.content = `${conditions.at(-1).varname()}`;
    for (let i = conditions.length - 2; i >= 0; i--) {
      const c = conditions[i]!;
      next.content = `(and ${c.varname()} ${next.content})`;
    }
    return next;
  }
  or(...conditions: ComplexCondition[]): ConditionBuilder {
    const next = new ConditionBuilder();
    // @ts-expect-error
    next.content = `${conditions.at(-1).varname()}`;
    for (let i = conditions.length - 2; i >= 0; i--) {
      const c = conditions[i]!;
      next.content = `(or ${c.varname()} ${next.content})`;
    }
    return next;
  }
  /** Equals to */
  equal(a: ValueProxy<any>, b: ValueProxy<any>): ConditionBuilder {
    const next = new ConditionBuilder();
    next.content = `(eq ${a.varname()} ${b.varname()})`;
    return next;
  }
  /** Not equals to */
  notEqual(a: ValueProxy<any>, b: ValueProxy<any>): ConditionBuilder {
    const next = new ConditionBuilder();
    next.content = `(ne ${a.varname()} ${b.varname()})`;
    return next;
  }
  /** Less than */
  lessThan(
    a: ValueProxy<Comparable>,
    b: ValueProxy<Comparable>
  ): ConditionBuilder {
    const next = new ConditionBuilder();
    next.content = `(lt ${a.varname()} ${b.varname()})`;
    return next;
  }
  /** Less than or equals to */
  lessEqualThan(
    a: ValueProxy<Comparable>,
    b: ValueProxy<Comparable>
  ): ConditionBuilder {
    const next = new ConditionBuilder();
    next.content = `(le ${a.varname()} ${b.varname()})`;
    return next;
  }
  /** Greater than */
  greaterThan(
    a: ValueProxy<Comparable>,
    b: ValueProxy<Comparable>
  ): ConditionBuilder {
    const next = new ConditionBuilder();
    next.content = `(gt ${a.varname()} ${b.varname()})`;
    return next;
  }
  /** Greater than or equals to */
  greaterEqualThan(
    a: ValueProxy<Comparable>,
    b: ValueProxy<Comparable>
  ): ConditionBuilder {
    const next = new ConditionBuilder();
    next.content = `(ge ${a.varname()} ${b.varname()})`;
    return next;
  }
  value(value: string | number | boolean): ValueProxy<any> {
    const varname = JSON.stringify(value);
    const toString = () => `{{${varname}}}`;
    return {
      varname: () => {
        return varname;
      },
      toString: toString,
      [Symbol.toHtmlTag]: toString,
      [Symbol.toPrimitive]: toString,
    };
  }
}
