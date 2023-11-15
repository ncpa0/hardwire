package htmxframework

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo"
	servestatic "github.com/ncpa0/htmx-framework/serve-static"
	templatebuilder "github.com/ncpa0/htmx-framework/template-builder"
	"github.com/ncpa0/htmx-framework/views"
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
	StaticURL        string
	BeforeStaticSend func(resp *servestatic.StaticResponse, c echo.Context) error
}

var conf *Configuration = &Configuration{
	StripExtension: false,
	DebugMode:      false,
	Entrypoint:     "index.tsx",
	ViewsDir:       "views",
	StaticDir:      "static",
	StaticURL:      "/static",
}

var DynamicResourceProvider *DRProvider = &DRProvider{
	resources: make(map[string]*DynamicResource),
}

func Configure(newConfig *Configuration) {
	templatebuilder.DebugMode = newConfig.DebugMode

	conf.StripExtension = newConfig.StripExtension
	conf.DebugMode = newConfig.DebugMode

	if newConfig.Entrypoint != "" {
		conf.Entrypoint = newConfig.Entrypoint
	}
	if newConfig.ViewsDir != "" {
		conf.ViewsDir = newConfig.ViewsDir
	}
	if newConfig.StaticDir != "" {
		conf.StaticDir = newConfig.StaticDir
	}
	if newConfig.StaticURL != "" {
		conf.StaticURL = newConfig.StaticURL
	}
	if newConfig.BeforeStaticSend != nil {
		conf.BeforeStaticSend = newConfig.BeforeStaticSend
	}
}

func createRouteHandler(view *views.View) func(c echo.Context) error {
	if conf.DebugMode {
		return func(c echo.Context) error {
			fmt.Printf("Received request for %s\n", c.Request().URL)
			selector := c.Request().Header.Get("HX-Target")

			c.Response().Header().Set("Cache-Control", "public, no-cache")
			ifNoneMatch := c.Request().Header.Get("If-None-Match")

			if selector != "" {
				fmt.Printf("  Applying html selector: %s\n", selector)
				child := view.QuerySelector("#" + selector)

				if !child.IsNil() {
					etag := child.Get().GetEtag()
					if ifNoneMatch == etag {
						fmt.Println("  Returning 304 Not Modified")
						return c.NoContent(http.StatusNotModified)
					}

					c.Response().Header().Set("ETag", etag)
					return c.HTML(http.StatusOK, child.Get().ToHtml())
				}
			}

			etag := view.GetNode().GetEtag()
			if ifNoneMatch == etag {
				fmt.Println("  Returning 304 Not Modified")
				return c.NoContent(http.StatusNotModified)
			}

			c.Response().Header().Set("ETag", etag)
			return c.HTML(http.StatusOK, view.GetNode().ToHtml())
		}
	}

	return func(c echo.Context) error {
		selector := c.Request().Header.Get("HX-Target")

		c.Response().Header().Set("Cache-Control", "public, no-cache")
		ifNoneMatch := c.Request().Header.Get("If-None-Match")

		if selector != "" {
			child := view.QuerySelector("#" + selector)

			if !child.IsNil() {
				if ifNoneMatch == child.Get().GetEtag() {
					return c.NoContent(http.StatusNotModified)
				}

				c.Response().Header().Set("ETag", child.Get().GetEtag())
				return c.HTML(http.StatusOK, child.Get().ToHtml())
			}
		}

		if ifNoneMatch == view.GetNode().GetEtag() {
			return c.NoContent(http.StatusNotModified)
		}

		c.Response().Header().Set("ETag", view.GetNode().GetEtag())
		return c.HTML(http.StatusOK, view.GetNode().ToHtml())
	}
}

func createDynamicFragmentHandler(view *views.View) func(c echo.Context) error {
	return func(c echo.Context) error {
		provider := DynamicResourceProvider.find(view.RequiredResource)
		resource, err := provider.handler(c)
		if err != nil {
			return err
		}
		if resource == nil {
			return c.NoContent(http.StatusNotFound)
		}

		var buff bytes.Buffer
		err = view.Template.Execute(&buff, resource)
		if err != nil {
			return err
		}
		return c.String(http.StatusOK, buff.String())
	}
}

func Start(server *echo.Echo, port string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	err = loadViews(wd)
	if err != nil {
		return err
	}

	addViewEndpoint := func(view *views.View) {
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
		server.GET(routePath, createRouteHandler(view))
	}

	addDynamicFragmentEndpoint := func(view *views.View) {
		var routePath string = view.GetFilepath()
		if !path.IsAbs(routePath) {
			routePath = "/" + view.GetFilepath()
		}
		routePath = routePath[:len(routePath)-len(".template.html")]

		if conf.DebugMode {
			fmt.Printf("Adding new dynamic fragment under route: %s\n", routePath)
		}
		server.GET(routePath, createDynamicFragmentHandler(view))
	}

	viewRegistry.ForEach(func(view *views.View) {
		if view.IsDynamicFragment {
			addDynamicFragmentEndpoint(view)
		} else {
			addViewEndpoint(view)
		}
	})

	if conf.DebugMode {
		fmt.Printf("Serving static files at the following URL: %s from directory: %s\n", conf.StaticURL, conf.StaticDir)
	}

	staticDir := conf.StaticDir
	if !path.IsAbs(staticDir) {
		staticDir = path.Join(wd, staticDir)
	}
	servestatic.Serve(server, conf.StaticURL, staticDir, &servestatic.Configuration{
		BeforeSend: conf.BeforeStaticSend,
	})

	return server.Start(port)
}
