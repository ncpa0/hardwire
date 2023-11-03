import { ComponentApi, defineContext } from "jsxte";
import { ValueProxy } from "./generate-go-templ";

type IfProps = JSXTE.PropsWithChildren<{
  condition: { varname(): string };
  negate?: boolean;
}>;

const ifCtx = defineContext<{ onElse(templ: string): void }>();

export const If = async (props: IfProps, comApi: ComponentApi) => {
  let elseTempl: string = "";
  const onElse = (templ: string) => {
    elseTempl = ` {{else}} ${templ} `;
  };

  const forTrue = await comApi.renderAsync(
    <ifCtx.Provider value={{ onElse }}>{props.children}</ifCtx.Provider>
  );

  return `{{if${
    props.negate ? " not " : ""
  } ${props.condition.varname()}}} ${forTrue}${elseTempl} {{end}}`;
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
  private content: string = "";

  private varname() {
    return this.content;
  }

  not(condition: ComplexCondition): ConditionBuilder {
    const next = new ConditionBuilder();
    // @ts-expect-error
    next.content = `not ${condition.varname()}`;
    return next;
  }
  and(...conditions: ComplexCondition[]): ConditionBuilder {
    const next = new ConditionBuilder();
    // @ts-expect-error
    next.content = `${conditions.at(-1).varname()}`;
    for (let i = conditions.length - 2; i >= 0; i--) {
      const c = conditions[i]!;
      // @ts-expect-error
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
      // @ts-expect-error
      next.content = `(or ${c.varname()} ${next.content})`;
    }
    return next;
  }
  /** Equals to */
  eq(a: ValueProxy<any>, b: ValueProxy<any>): ConditionBuilder {
    const next = new ConditionBuilder();
    next.content = `(eq ${a.varname()} ${b.varname()})`;
    return next;
  }
  /** Not equals to */
  ne(a: ValueProxy<any>, b: ValueProxy<any>): ConditionBuilder {
    const next = new ConditionBuilder();
    next.content = `(ne ${a.varname()} ${b.varname()})`;
    return next;
  }
  /** Less than */
  lt(a: ValueProxy<Comparable>, b: ValueProxy<Comparable>): ConditionBuilder {
    const next = new ConditionBuilder();
    next.content = `(lt ${a.varname()} ${b.varname()})`;
    return next;
  }
  /** Less than or equals to */
  le(a: ValueProxy<Comparable>, b: ValueProxy<Comparable>): ConditionBuilder {
    const next = new ConditionBuilder();
    next.content = `(le ${a.varname()} ${b.varname()})`;
    return next;
  }
  /** Greater than */
  gt(a: ValueProxy<Comparable>, b: ValueProxy<Comparable>): ConditionBuilder {
    const next = new ConditionBuilder();
    next.content = `(gt ${a.varname()} ${b.varname()})`;
    return next;
  }
  /** Greater than or equals to */
  ge(a: ValueProxy<Comparable>, b: ValueProxy<Comparable>): ConditionBuilder {
    const next = new ConditionBuilder();
    next.content = `(ge ${a.varname()} ${b.varname()})`;
    return next;
  }

  value(value: string | number | boolean): ValueProxy<any> {
    return {
      varname: () => {
        return JSON.stringify(value);
      },
      toString: () => {
        return `{{${JSON.stringify(value)}}}`;
      },
    };
  }
}

export const condition = (
  factory: (builder: ConditionBuilder) => ConditionBuilder
) => {
  const builder = new ConditionBuilder();
  return factory(builder) as any as { varname(): string };
};
