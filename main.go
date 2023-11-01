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

type Configuration struct {
	// When enabled, the `.html` extension will be stripped from the URL pathnames.
	StripExtension bool
	// When enabled, the server will print debug information to the console.
	DebugMode bool
	// The entrypoint file containing the JSX pages used to generate the views html files.
	//
	// Defaults to `index.tsx`.
	Entrypoint string
	// The directory to which output the generated html files, and from which those will be hosted.
	//
	// Defaults to `views`.
	ViewsDir string
	// The directory to which output the static files, and from which those will be hosted.
	//
	// Defaults to `static`.
	StaticDir string
	// The URL path from under which the static files will be hosted.
	//
	// Defaults to `/static`.
	StaticURL string
}

var viewRegistry = views.NewViewRegistry()

var conf *Configuration = &Configuration{
	StripExtension: true,
	DebugMode:      false,
	Entrypoint:     "index.tsx",
	ViewsDir:       "views",
	StaticDir:      "static",
	StaticURL:      "/static",
}

func Configure(config *Configuration) {
	conf = config
	templatebuilder.DebugMode = config.DebugMode
}

func loadViews() error {
	wd, _ := os.Getwd()
	viewsFullPath := conf.ViewsDir
	if !path.IsAbs(viewsFullPath) {
		viewsFullPath = path.Join(wd, viewsFullPath)
	}

	err := templatebuilder.BuildPages(conf.Entrypoint, viewsFullPath, conf.StaticDir, conf.StaticURL)
	if err != nil {
		return err
	}

	if conf.DebugMode {
		fmt.Printf("Loading view from %s\n", viewsFullPath)
	}
	err = utils.Walk(viewsFullPath, func(root string, dirs []string, files []string) error {
		for _, file := range files {
			ext := path.Ext(file)

			if ext != ".html" {
				continue
			}

			fullPath := path.Join(root, file)
			relToView := fullPath[len(viewsFullPath):]
			view, err := views.NewView(viewsFullPath, relToView)
			if err != nil {
				return err
			}

			if conf.DebugMode {
				fmt.Printf("Loading view from file %s\n", file)
				fmt.Printf("ROOT: %s PATH: %s\n", viewsFullPath, relToView)
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

func createRouteHandler(view *views.View) func(c echo.Context) error {
	if conf.DebugMode {
		return func(c echo.Context) error {
			fmt.Printf("Received request for %s\n", c.Request().URL)
			selector := c.Request().Header.Get("HX-Target")

			if selector != "" {
				fmt.Printf("Applying html selector: %s\n", selector)
				child := view.QuerySelector("#" + selector)

				if !child.IsNil() {
					err := renderNode(c, child.Get())
					if err != nil {
						fmt.Println("Rendering the node has failed.")
					}
					return err
				}
			}

			err := renderNode(c, view.ToNode())
			if err != nil {
				fmt.Println("Rendering the node has failed.")
			}
			return err
		}
	}

	return func(c echo.Context) error {
		selector := c.Request().Header.Get("HX-Target")

		if selector != "" {
			child := view.QuerySelector("#" + selector)

			if !child.IsNil() {
				return renderNode(c, child.Get())
			}
		}

		return renderNode(c, view.ToNode())
	}
}

func Start(e *echo.Echo, port string) error {
	err := loadViews()
	if err != nil {
		return err
	}

	viewRegistry.ForEach(func(view *views.View) {
		var routePath string = view.GetFilepath()
		if !path.IsAbs(routePath) {
			routePath = "/" + view.GetFilepath()
		}
		if conf.StripExtension {
			routePath = routePath[:len(routePath)-len(path.Ext(routePath))]
		}

		if conf.DebugMode {
			fmt.Printf("Adding new route: %s\n", routePath)
		}
		e.GET(routePath, createRouteHandler(view))
	})

	if conf.DebugMode {
		fmt.Printf("Serving static files at the following URL: %s from directory: %s\n", conf.StaticURL, conf.StaticDir)
	}
	e.Static(conf.StaticURL, conf.StaticDir)

	return e.Start(port)
}
