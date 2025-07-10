import fs from "node:fs";
import path from "node:path";
import { ftmpl } from "../utils/ftmpl";

const BUNFIG_TEMPLATE = ftmpl`
jsx = "react-jsx"\njsxImportSource = "jsxte"
`;
const TSCONFIG_TEMPLATE = JSON.stringify(
  {
    compilerOptions: {
      target: "ES2022",
      module: "ES2022",
      moduleResolution: "Node",
      jsx: "react-jsx",
      jsxImportSource: "jsxte",
    },
  },
  null,
  2,
);
const CSS_TEMPLATE = ftmpl`
body {
  margin: unset;
  min-height: 100vh;
}
`;
const INDEX_TEMPLATE = ftmpl`
import "hardwire-html-generator";

export default function App() {
  return (
    <Html nochunked>
      <Head>
        <Style dirname={import.meta.dir} path="./style.css" />
        <Script dirname={import.meta.dir} path="./index.client.ts" />
      </Head>
      <div id="root">
        <nav>
          <ul>
            <li><Link href="/">Home</Link></li>
          </ul>
        </nav>
        <Switch id="main-switch">
          <StaticRoute path="home" title="Home Page">
            <h1>Hello World</h1>
          </StaticRoute>
        </Switch>
      </div>
    </html>
  );
}`;

export async function initCmd(wd: string) {
  const bunfigPath = path.join(wd, "bunfig.toml");
  if (!fs.existsSync(bunfigPath)) {
    await Bun.write(bunfigPath, BUNFIG_TEMPLATE);
  }

  const indexPath = path.join(wd, "index.tsx");
  if (!fs.existsSync(indexPath)) {
    await Bun.write(indexPath, INDEX_TEMPLATE);

    const tsconfigPath = path.join(wd, "tsconfig.json");
    if (!fs.existsSync(tsconfigPath)) {
      await Bun.write(tsconfigPath, TSCONFIG_TEMPLATE);
    }

    const stylePath = path.join(wd, "style.css");
    if (!fs.existsSync(stylePath)) {
      await Bun.write(stylePath, CSS_TEMPLATE);
    }

    const indexClientPath = path.join(wd, "index.client.ts");
    if (!fs.existsSync(indexClientPath)) {
      await Bun.write(indexClientPath, "");
    }
  }
}
