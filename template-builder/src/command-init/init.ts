import fs from "node:fs/promises";
import path from "node:path";

const BUNFIG_TEMPLATE = `jsx = "react-jsx"\njsxImportSource = "jsxte"\n`;
const TSCONFIG_TEMPLATE = `{
    "compilerOptions": {
        "jsx": "react-jsx",
        "jsxImportSource": "jsxte"
    }
}
`;
const CSS_TEMPLATE = `body { 
    margin: unset;
    min-height: 100vh;
    min-width: 100vw;
  }
  
  #root {
    display: flex;
    flex-direction: column;
  }`;
const INDEX_TEMPLATE = `import "hardwire-html-generator";
        
export default function App() {
  return (
    <html>
      <Head>
        <meta charset='utf-8' />
        <meta http-equiv='X-UA-Compatible' content='IE=edge' />
        <meta name='viewport' content='width=device-width, initial-scale=1' />
        <Style path="./style.css" dirname={import.meta.dir} />
        <Script path="./index.client.ts" dirname={import.meta.dir} />
      </Head>
      <body>
        <div id="root">
          <nav>
            <ul>
              <li><Link href="/">Home</Link></li>
            </ul>
          </nav>
          <Switch>
            <StaticRoute path="home" title="Home Page">
              <h1>Hello World</h1>
            </StaticRoute>
          </Switch>
        </div>
      </body>
    </html>
  );
}`;

export async function initCmd(wd: string) {
  const bunfigPath = path.join(wd, "bunfig.toml");
  if (!(await fs.exists(bunfigPath))) {
    await Bun.write(bunfigPath, BUNFIG_TEMPLATE);
  }

  const indexPath = path.join(wd, "index.tsx");
  if (!(await fs.exists(indexPath))) {
    await Bun.write(indexPath, INDEX_TEMPLATE);

    const tsconfigPath = path.join(wd, "tsconfig.json");
    if (!(await fs.exists(tsconfigPath))) {
      await Bun.write(tsconfigPath, TSCONFIG_TEMPLATE);
    }

    const stylePath = path.join(wd, "style.css");
    if (!(await fs.exists(stylePath))) {
      await Bun.write(stylePath, CSS_TEMPLATE);
    }

    const indexClientPath = path.join(wd, "index.client.ts");
    if (!(await fs.exists(indexClientPath))) {
      await Bun.write(indexClientPath, "");
    }
  }
}
