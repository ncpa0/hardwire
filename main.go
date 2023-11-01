package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo"
	"github.com/ncpa0cpl/blazing/utils"
	"github.com/ncpa0cpl/blazing/views"
	"golang.org/x/net/html"
)

var VIEWS = utils.GetEnv("VIEWS_DIR", "./")
var STATIC_DIR = utils.GetEnv("STATIC_DIR", "./static")
var STATIC_URL = utils.GetEnv("STATIC_URL", "/static")
var ViewRegistry = views.NewViewRegistry()

func loadViews() {
	wd, _ := os.Getwd()
	viewsFullPath := path.Join(wd, VIEWS)

	fmt.Println(viewsFullPath)

	err := utils.Walk(viewsFullPath, func(root string, dirs []string, files []string) error {
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
			ViewRegistry.Register(view)
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error loading views.")
		fmt.Println(err)
	}
}

func RenderView(c echo.Context, viewKey string) error {
	view := ViewRegistry.GetView("/index.html")

	if view.IsNil() {
		return c.String(http.StatusNotFound, "Not Found")
	}

	node := view.Get().ToNode()

	var b bytes.Buffer
	err := html.Render(&b, node)

	if err != nil {
		return err
	}

	return c.HTML(http.StatusOK, b.String())
}

func RenderNode(c echo.Context, node *html.Node) error {
	var b bytes.Buffer
	err := html.Render(&b, node)

	if err != nil {
		return err
	}

	return c.HTML(http.StatusOK, b.String())
}

func main() {
	loadViews()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, "/index.html")
	})

	ViewRegistry.ForEach(func(view *views.View) {
		fmt.Println("Adding route: ", view.GetFilepath())
		e.GET(view.GetFilepath(), func(c echo.Context) error {
			selector := c.Request().Header.Get("HX-Target")

			if selector != "" {
				child := view.QuerySelector("#" + selector)

				if !child.IsNil() {
					return RenderNode(c, child.Get())
				}
			}

			return RenderNode(c, view.ToNode())
		})
	})

	e.Static(STATIC_URL, STATIC_DIR)

	e.Logger.Fatal(e.Start(":8080"))
}
