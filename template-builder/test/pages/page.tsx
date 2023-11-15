import { StreamedFragment } from "./streamed-fragment";

export default function MainPage() {
  return (
    <html>
      <head>
        <title>Test</title>
        <Script dirname={import.meta.dir} path="./client.ts" />
        <Script package="adwaveui" />
        <Script embed>console.log("embedded script");</Script>
        <Style dirname={import.meta.dir} path="./style.css" />
        <Style package="adwavecss/dist/styles.css" />
        <Style embed>
          {`.body {
            background-color: red;
          }`}
        </Style>
      </head>
      <body>
        <nav>
          <ul>
            <li>
              <Link href="/home">Home</Link>
            </li>
            <li>
              <Link href="/about">Home</Link>
            </li>
            <li>
              <Link href="/products">Home</Link>
            </li>
            <li>
              <Link href="/products/1">Home</Link>
            </li>
          </ul>
        </nav>
        <Router id="root">
          <Route path="home">
            <h1>Home</h1>
            <Link href="/home/sub1/sub2" />
            <Router id="s1">
              <Route path="sub1">
                <h2>Sub1</h2>
                <Link href="/home/sub1/sub2" />
                <Router id="s2">
                  <Route path="sub2">
                    <h3>Sub2</h3>
                  </Route>
                </Router>
              </Route>
            </Router>
          </Route>
          <Route path="about">
            <h1>About</h1>
          </Route>
          <Route path="products">
            <h2>Products</h2>
            <Link href="/products/1">Open 1</Link>
            <Link href="/products/2">Open 2</Link>
            <Link href="/products/3">Open 3</Link>
            <Router id="products">
              <Route path="1">
                <span> Product 1 </span>
              </Route>
              <Route path="2">
                <span> Product 2 </span>
              </Route>
              <Route path="3">
                <span> Product 3 </span>
              </Route>
            </Router>
          </Route>
          <Route path="user">
            <StreamedFragment />
          </Route>
        </Router>
      </body>
    </html>
  );
}
