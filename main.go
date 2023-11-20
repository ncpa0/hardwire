package htmxframework

import (
	"fmt"
	"os"
	"path"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo"
	config "github.com/ncpa0/htmx-framework/configuration"
	resources "github.com/ncpa0/htmx-framework/resource-provider"
	servestatic "github.com/ncpa0/htmx-framework/serve-static"
	"github.com/ncpa0/htmx-framework/views"
)

type DynamicRequestContext = resources.DynamicRequestContext
type ResourceRequestError = resources.ResourceRequestError
type DRProvider = resources.DRProvider
type Configuration = config.Configuration
type CachingConfig = config.CachingConfig
type CachingPolicy = config.CachingPolicy

var ResourceProvider = resources.Provider
var Configure = config.Configure

func Start(server *echo.Echo, port string) error {
	pageViewRegistry := views.GetPageViewRegistry()
	dynamicFragmentViewRegistry := views.GetDynamicFragmentViewRegistry()

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	err = views.LoadViews(wd)
	if err != nil {
		return err
	}

	pageViewRegistry.ForEach(func(view *views.PageView) {
		if config.Current.DebugMode {
			fmt.Printf("Adding new route: %s\n", view.GetRoutePathname())
		}
		server.GET(view.GetRoutePathname(), createPageViewHandler(view))
	})

	dynamicFragmentViewRegistry.ForEach(func(view *views.DynamicFragmentView) {
		if config.Current.DebugMode {
			fmt.Printf("Adding new dynamic fragment under route: %s\n", view.GetRoutePathname())
		}
		server.GET(view.GetRoutePathname(), createDynamicFragmentHandler(view))
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
