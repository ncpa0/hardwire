package hardwire

import (
	"net/http"

	echo "github.com/labstack/echo/v4"
	config "github.com/ncpa0/hardwire/configuration"
	resources "github.com/ncpa0/hardwire/resources"
	"github.com/ncpa0/hardwire/utils"
	"github.com/ncpa0/hardwire/views"
)

func createDynamicFragmentHandler(view *views.DynamicFragmentView, conf *config.Configuration) func(c echo.Context) error {
	return func(c echo.Context) error {
		resKey := view.ResourceKeys()[0]
		isValidResource := resources.HasResource(resKey)

		if !isValidResource {
			err := c.NoContent(http.StatusNotFound)
			if err != nil {
				return err
			}
			if conf.BeforeResponse != nil {
				return conf.BeforeResponse(c)
			}
			return nil
		}

		routePathname := c.Request().Header.Get("Hardwire-Dynamic-Fragment-Request")
		if routePathname == "" {
			err := c.String(http.StatusBadRequest, "Bad Request")
			if err != nil {
				return err
			}
			if conf.BeforeResponse != nil {
				return conf.BeforeResponse(c)
			}
			return nil
		}

		hxCurrentUrl := c.Request().Header.Get("Hx-Current-Url")
		params := utils.ParseUrlParams(routePathname, hxCurrentUrl)

		handler, err := HardwireContext.GetResourceHandler(c, resKey)
		if err != nil {
			return err
		}

		resource, err := handler(routePathname, params)

		if err != nil {
			err := utils.HandleError(c, err)
			if err != nil {
				return err
			}
			if conf.BeforeResponse != nil {
				return conf.BeforeResponse(c)
			}
			return nil
		}
		if resource == nil {
			err := c.String(http.StatusNotFound, "Not found")
			if err != nil {
				return err
			}
			if conf.BeforeResponse != nil {
				return conf.BeforeResponse(c)
			}
			return nil
		}

		html, err := view.Build(resource)
		if err != nil {
			err := utils.HandleError(c, err)
			if err != nil {
				return err
			}
			if conf.BeforeResponse != nil {
				return conf.BeforeResponse(c)
			}
			return nil
		}

		c.Response().Header().Set(
			"Vary",
			"Hx-Current-Url, Hardwire-Dynamic-Fragment-Request, Accept-Language",
		)

		if !config.Current.Caching.Fragments.NoStore {
			etag := utils.Hash(html)
			ifNoneMatch := c.Request().Header.Get("If-None-Match")
			if ifNoneMatch == etag {
				err := c.NoContent(http.StatusNotModified)
				if err != nil {
					return err
				}
				if conf.BeforeResponse != nil {
					return conf.BeforeResponse(c)
				}
				return nil
			}

			c.Response().Header().Set("ETag", etag)
		}

		c.Response().Header().Set(
			"Cache-Control",
			config.GenerateCacheHeaderForFragments(),
		)
		c.Response().Header().Set("Content-Type", "text/html")
		err = c.String(http.StatusOK, html)
		if err != nil {
			return err
		}
		if conf.BeforeResponse != nil {
			return conf.BeforeResponse(c)
		}
		return nil
	}
}
