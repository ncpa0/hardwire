## Installation

```sh
go get github.com/ncpa0/hardwire
```

## Setup

```go
import (
    "github.com/ncpa0/hardwire"
    echo "github.com/labstack/echo/v4"
)

func main() {
    server := echo.New()

    hardwire.Configure(&hardwire.Configuration{
        Entrypoint: "./src/index.tsx",
        ViewsDir: "./pages",
        StaticDir: "./static",
        StaticURL: "/static",
    })

    server.GET("/", func(c echo.Context) error {
		return c.Redirect(303, "/home")
	})

    server.Logger.Fatal(hardwire.Start(server, ":8080"))
}
```
