package htmxframework

import (
	"fmt"
	"net/http"
	"os"
	"path"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo"
	config "github.com/ncpa0/htmx-framework/configuration"
	servestatic "github.com/ncpa0/htmx-framework/serve-static"
	"github.com/ncpa0/htmx-framework/views"
)

var DynamicResourceProvider *DRProvider = &DRProvider{
	resources: make(map[string]*DynamicResource),
}

func Configure(newConfig *config.Configuration) {
	config.Current.StripExtension = newConfig.StripExtension
	config.Current.DebugMode = newConfig.DebugMode

	if newConfig.Entrypoint != "" {
		config.Current.Entrypoint = newConfig.Entrypoint
	}
	if newConfig.ViewsDir != "" {
		config.Current.ViewsDir = newConfig.ViewsDir
	}
	if newConfig.StaticDir != "" {
		config.Current.StaticDir = newConfig.StaticDir
	}
	if newConfig.StaticURL != "" {
		config.Current.StaticURL = newConfig.StaticURL
	}
	if newConfig.BeforeStaticSend != nil {
		config.Current.BeforeStaticSend = newConfig.BeforeStaticSend
	}
}

func createRouteHandler(view *views.View) func(c echo.Context) error {
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
				c.Response().Header().Set("Content-Type", "text/html")
				return c.String(http.StatusOK, child.Get().ToHtml())
			}
		}

		if ifNoneMatch == view.GetNode().GetEtag() {
			return c.NoContent(http.StatusNotModified)
		}

		c.Response().Header().Set("ETag", view.GetNode().GetEtag())
		c.Response().Header().Set("Content-Type", "text/html")
		return c.String(http.StatusOK, view.GetNode().ToHtml())
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
		if config.Current.DebugMode {
			fmt.Printf("Adding new route: %s\n", view.GetRoutePathname())
		}
		server.GET(view.GetRoutePathname(), createRouteHandler(view))
	}

	addDynamicFragmentEndpoint := func(view *views.View) {
		if config.Current.DebugMode {
			fmt.Printf("Adding new dynamic fragment under route: %s\n", view.GetRoutePathname())
		}
		server.GET(view.GetRoutePathname(), createDynamicFragmentHandler(view))
	}

	viewRegistry.ForEach(func(view *views.View) {
		if view.IsDynamicFragment {
			addDynamicFragmentEndpoint(view)
		} else {
			addViewEndpoint(view)
		}
	})

	registerActionHandlers(server)

	if config.Current.DebugMode {
		fmt.Printf(
			"Serving static files at the following URL: %s from directory: %s\n",
			config.Current.StaticURL, config.Current.StaticDir,
		)
	}

	staticDir := config.Current.StaticDir
	if !path.IsAbs(staticDir) {
		staticDir = path.Join(wd, staticDir)
	}

	servestatic.Serve(server, config.Current.StaticURL, staticDir, &servestatic.Configuration{
		BeforeSend: config.Current.BeforeStaticSend,
	})

	return server.Start(port)
}
