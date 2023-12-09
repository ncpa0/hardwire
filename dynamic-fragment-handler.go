package hardwire

import (
	"net/http"

	"github.com/labstack/echo"
	config "github.com/ncpa0/hardwire/configuration"
	resources "github.com/ncpa0/hardwire/resource-provider"
	"github.com/ncpa0/hardwire/utils"
	"github.com/ncpa0/hardwire/views"
)

func createDynamicFragmentHandler(view *views.DynamicFragmentView) func(c echo.Context) error {
	return func(c echo.Context) error {
		provider := resources.Provider.Find(view.GetResourceName())

		if provider.IsNil() {
			return c.NoContent(http.StatusNotFound)
		}

		routePathname := c.Request().Header.Get("HX-Dynamic-Fragment-Request")
		if routePathname == "" {
			return c.String(http.StatusBadRequest, "Bad Request")
		}

		hxCurrentUrl := c.Request().Header.Get("Hx-Current-Url")
		params := utils.ParseUrlParams(routePathname, hxCurrentUrl)

		requestContext := resources.NewDynamicRequestContext(
			c,
			params,
			routePathname,
		)

		resource, err := provider.Get().Handle(requestContext)
		if err != nil {
			return utils.HandleError(c, err)
		}
		if resource == nil {
			return c.String(http.StatusNotFound, "Not found")
		}

		html, err := view.Build(resource)
		if err != nil {
			return utils.HandleError(c, err)
		}

		c.Response().Header().Set("Vary", "Hx-Current-Url, HX-Dynamic-Fragment-Request")

		if !config.Current.Caching.Fragments.NoStore {
			etag := utils.Hash(html)
			ifNoneMatch := c.Request().Header.Get("If-None-Match")
			if ifNoneMatch == etag {
				return c.NoContent(http.StatusNotModified)
			}

			c.Response().Header().Set("ETag", etag)
		}

		c.Response().Header().Set(
			"Cache-Control",
			config.GenerateCacheHeaderForFragments(),
		)
		c.Response().Header().Set("Content-Type", "text/html")
		return c.String(http.StatusOK, html)
	}
}
