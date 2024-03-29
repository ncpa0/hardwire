import { ComponentApi, ElementGenerator, JsxteRenderer } from "jsxte";

const SELF_CLOSING_TAG_LIST = [
  "area",
  "base",
  "br",
  "col",
  "embed",
  "hr",
  "img",
  "input",
  "link",
  "meta",
  "param",
  "source",
  "track",
  "wbr",
];

const sanitizeAttributeValue = (value: string): string => {
  return value
    .replaceAll('"', "&quot;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;");
};

const attributeToHtmlTagString = ([key, value]: [
  string,
  string | boolean | number | undefined,
]): string => {
  if (value === true) {
    return `${key}`;
  }
  if (value === false || value === null || value === undefined) {
    return "";
  }
  return `${key}="${sanitizeAttributeValue(String(value))}"`;
};

const mapAttributesToHtmlTagString = (attributes: [string, any][]): string => {
  const results: string[] = [];

  for (let i = 0; i < attributes.length; i++) {
    const attribute = attributes[i]!;
    const html = attributeToHtmlTagString(attribute);
    if (html.length > 0) results.push(html);
  }

  return results.join(" ");
};

class BaseHtmlGenerator {
  constructor() {}

  protected generateTag(
    tag: string,
    attributes?: string,
    content?: string
  ): string {
    if (attributes) {
      attributes = " " + attributes;
    } else {
      attributes = "";
    }

    if (!content || content.length === 0) {
      if (SELF_CLOSING_TAG_LIST.includes(tag)) {
        return `<${tag}${attributes} />`;
      } else {
        return `<${tag}${attributes}></${tag}>`;
      }
    }

    return `<${tag}${attributes}>${content}</${tag}>`;
  }

  protected flattenChildren(children: string[]): string {
    return children.join("");
  }

  public static leftPad(str: string, pad: string): string {
    if (!str.includes("\n")) {
      return pad + str;
    } else {
      const lines = str.split("\n");
      for (let i = 0; i < lines.length; i++) {
        lines[i] = pad + lines[i];
      }
      return lines.join("\n");
    }
  }

  public static trimContent(content: string): {
    wsLeft: boolean;
    wsRight: boolean;
    trimmed: string;
  } {
    let leftWhitespace = 0;
    let rightWhitespace = 0;
    let wsLeft = false;
    let wsRight = false;

    for (let i = 0; i < content.length; i++) {
      if (content[i] === " " || content[i] === "\n") {
        leftWhitespace += 1;
      } else {
        break;
      }
    }

    if (leftWhitespace === content.length) {
      return { wsLeft: true, wsRight: true, trimmed: "" };
    }

    if (leftWhitespace > 0) {
      content = content.substring(leftWhitespace);
      wsLeft = true;
    }

    for (let i = content.length - 1; i >= 0; i--) {
      if (content[i] === " " || content[i] === "\n") {
        rightWhitespace += 1;
      } else {
        break;
      }
    }
    if (rightWhitespace > 0) {
      content = content.substring(0, content.length - rightWhitespace);
      wsRight = true;
    }

    return { wsLeft, wsRight, trimmed: content };
  }
}

class HtmlCompactGenerator
  extends BaseHtmlGenerator
  implements ElementGenerator<string | Promise<string>>
{
  createTextNode(text: string | number | bigint): string {
    return String(text);
  }

  createElement(
    type: string,
    attributes: [attributeName: string, attributeValue: any][],
    children: string[]
  ): Promise<string> {
    return Promise.resolve(children)
      .then((c) => Promise.all(c))
      .then((children) => {
        const attributesString = mapAttributesToHtmlTagString(attributes);
        const content = this.flattenChildren(children);
        return this.generateTag(type, attributesString, content);
      });
  }

  createFragment(children: string[]): Promise<string> {
    return Promise.resolve(children)
      .then((c) => Promise.all(c))
      .then((children) => {
        return this.flattenChildren(children);
      });
  }
}

const generator = new HtmlCompactGenerator();

const renderer = new JsxteRenderer<string | Promise<string>>(generator, {
  allowAsync: true,
});

export const render = (element: JSX.Element, api?: ComponentApi) => {
  if (api) {
    const renderer = new JsxteRenderer<string | Promise<string>>(
      generator,
      {
        allowAsync: true,
      },
      api
    );
    return renderer.render(element);
  }
  return renderer.render(element);
};
