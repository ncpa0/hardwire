import { describe, expect, it } from "bun:test";
import path from "node:path";

const file = (fp: string) => Bun.file(path.resolve(import.meta.dir, fp));

describe("build", () => {
  it("correctly builds html files", async () => {
    if (await file("./tmp").exists()) {
      Bun.spawnSync({
        cwd: path.dirname(import.meta.dir),
        cmd: ["rm", "-rf", "./tmp"],
      });
    }

    const res = Bun.spawnSync({
      cwd: path.dirname(import.meta.dir),
      cmd: [
        "bun",
        "./src/index.tsx",
        "build",
        "--src",
        "./test/pages/page.tsx",
        "--outdir",
        "./test/tmp",
      ],
    });

    if (!res.success) {
      throw new Error(res.stderr.toString("utf-8"));
    }

    const home = file("./tmp/home.html");
    const about = file("./tmp/about.html");
    const products = file("./tmp/products.html");
    const homeSub1 = file("./tmp/home/sub1.html");
    const homeSub2 = file("./tmp/home/sub1/sub2.html");
    const product1 = file("./tmp/products/1.html");
    const product2 = file("./tmp/products/2.html");
    const product3 = file("./tmp/products/3.html");

    expect(await home.text()).toMatchSnapshot();
    expect(await about.text()).toMatchSnapshot();
    expect(await products.text()).toMatchSnapshot();
    expect(await homeSub1.text()).toMatchSnapshot();
    expect(await homeSub2.text()).toMatchSnapshot();
    expect(await product1.text()).toMatchSnapshot();
    expect(await product2.text()).toMatchSnapshot();
    expect(await product3.text()).toMatchSnapshot();
  });
});
