package hardwire

import (
	"fmt"
	"os"
	"path"

	"github.com/labstack/echo"
	config "github.com/ncpa0/hardwire/configuration"
	resources "github.com/ncpa0/hardwire/resource-provider"
	servestatic "github.com/ncpa0/hardwire/serve-static"
	"github.com/ncpa0/hardwire/views"
)

type DynamicRequestContext = resources.DynamicRequestContext
type ResourceRequestError = resources.ResourceRequestError
type DRProvider = resources.DRProvider
type Configuration = config.Configuration
type CachingConfig = config.CachingConfig
type CachingPolicy = config.CachingPolicy

var ResourceProvider = resources.Provider
var Configure = config.Configure

func redirectHandler(to string) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return ctx.Redirect(301, to)
	}
}

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
		pathname := view.GetRoutePathname()
		server.GET(pathname, createPageViewHandler(view))
		server.GET(pathname+"/", redirectHandler(pathname))
	})

	dynamicFragmentViewRegistry.ForEach(func(view *views.DynamicFragmentView) {
		if config.Current.DebugMode {
			fmt.Printf("Adding new dynamic fragment under route: %s\n", view.GetRoutePathname())
		}
		pathname := view.GetRoutePathname()
		server.GET(pathname, createDynamicFragmentHandler(view))
		server.GET(pathname+"/", redirectHandler(pathname))
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
