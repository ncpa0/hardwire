package htmxframework

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo"
	templatebuilder "github.com/ncpa0/htmx-framework/template-builder"
	"github.com/ncpa0/htmx-framework/utils"
	"github.com/ncpa0/htmx-framework/views"
	"golang.org/x/net/html"
)

var ENTRYPOINT = utils.GetEnv("PAGE_ENTRYPOINT", "./pages/index.tsx")
var VIEWS = utils.GetEnv("VIEWS_DIR", "./")
var STATIC_DIR = utils.GetEnv("STATIC_DIR", "./static")
var STATIC_URL = utils.GetEnv("STATIC_URL", "/static")
var viewRegistry = views.NewViewRegistry()

func SetDebugMode(debugMode bool) {
	templatebuilder.DebugMode = debugMode
}

func loadViews() error {
	wd, _ := os.Getwd()
	viewsFullPath := VIEWS
	if !path.IsAbs(viewsFullPath) {
		viewsFullPath = path.Join(wd, viewsFullPath)
	}

	err := templatebuilder.BuildPages(ENTRYPOINT, viewsFullPath)
	if err != nil {
		return err
	}

	err = utils.Walk(viewsFullPath, func(root string, dirs []string, files []string) error {
		for _, file := range files {
			ext := path.Ext(file)

			if ext != ".html" {
				continue
			}

			fullPath := path.Join(root, file)
			relToView := fullPath[len(viewsFullPath):]
			view, err := views.NewView(VIEWS, relToView)
			if err != nil {
				return err
			}
			viewRegistry.Register(view)
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error loading views.")
		return err
	}

	return nil
}

func renderNode(c echo.Context, node *html.Node) error {
	var b bytes.Buffer
	err := html.Render(&b, node)

	if err != nil {
		return err
	}

	return c.HTML(http.StatusOK, b.String())
}

func Start(e *echo.Echo, port string) error {
	err := loadViews()
	if err != nil {
		return err
	}

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, "/index.html")
	})

	viewRegistry.ForEach(func(view *views.View) {
		e.GET(view.GetFilepath(), func(c echo.Context) error {
			selector := c.Request().Header.Get("HX-Target")

			if selector != "" {
				child := view.QuerySelector("#" + selector)

				if !child.IsNil() {
					return renderNode(c, child.Get())
				}
			}

			return renderNode(c, view.ToNode())
		})
	})

	e.Static(STATIC_URL, STATIC_DIR)

	return e.Start(port)
}
