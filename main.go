package hardwire

import (
	"fmt"
	"os"
	"path"

	echo "github.com/labstack/echo/v4"
	config "github.com/ncpa0/hardwire/configuration"
	hw "github.com/ncpa0/hardwire/hw-context"
	resources "github.com/ncpa0/hardwire/resources"
	servestatic "github.com/ncpa0/hardwire/serve-static"
	"github.com/ncpa0/hardwire/views"
)

type DynamicRequestContext = resources.DynamicRequestContext
type ResourceRequestError = resources.ResourceRequestError
type ResourceRegistry = resources.ResourceRegistry
type ActionContext = resources.ActionContext
type Configuration = config.Configuration
type CachingConfig = config.CachingConfig
type CachingPolicy = config.CachingPolicy

var ResourceReg = resources.ResourceReg
var Configure = config.Configure
var HardwireContext hw.HardwireContext = &HwContext{}

func redirectHandler(to string) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return ctx.Redirect(301, to)
	}
}

func validateResourcesAvailable(resource []string) error {
	for _, resKey := range resource {
		_, found := resources.GetResource(resKey)

		if !found {
			return fmt.Errorf("'%s' resource doesn't have a provider registered", resKey)
		}
	}

	return nil
}

// Builds the HTML and templates for all pages and adds the routes to the server
func UseWith(server *echo.Echo) error {
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

	resources.ValidateActionEndpoints()

	err = pageViewRegistry.ForEach(func(view *views.PageView) error {
		fmt.Printf("Adding new route: %s\n", view.GetRoutePathname())

		if view.IsDynamic() {
			err := validateResourcesAvailable(view.GetResourceKeys().ToSlice())
			if err != nil {
				return err
			}
		}

		pathname := view.GetRoutePathname()
		server.GET(pathname, createPageViewHandler(view, config.Current))
		server.GET(pathname+"/", redirectHandler(pathname))

		return nil
	})

	if err != nil {
		return err
	}

	err = dynamicFragmentViewRegistry.ForEach(func(view *views.DynamicFragmentView) error {
		if config.Current.DebugMode {
			fmt.Printf("Adding new dynamic fragment under route: %s\n", view.GetRoutePathname())
		}

		resKey := view.ResourceKeys()[0]
		err := validateResourcesAvailable([]string{resKey})
		if err != nil {
			return err
		}

		pathname := view.GetRoutePathname()
		server.GET(pathname, createDynamicFragmentHandler(view, config.Current))
		server.GET(pathname+"/", redirectHandler(pathname))

		return nil
	})

	if err != nil {
		return err
	}

	resources.MountActionEndpoints(HardwireContext, server)

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
		BeforeSend: config.Current.BeforeStaticResponse,
	})

	return nil
}
