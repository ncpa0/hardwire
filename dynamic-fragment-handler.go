package hardwire

import (
	"errors"
	"net/http"

	echo "github.com/labstack/echo/v4"
	config "github.com/ncpa0/hardwire/configuration"
	resources "github.com/ncpa0/hardwire/resource-provider"
	"github.com/ncpa0/hardwire/utils"
	"github.com/ncpa0/hardwire/views"
)

func createDynamicFragmentHandler(view *views.DynamicFragmentView, conf *config.Configuration) func(c echo.Context) error {
	return func(c echo.Context) error {
		provider := resources.Provider.Find(view.GetResourceName())

		if provider.IsNil() {
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

		requestContext := resources.NewDynamicRequestContext(
			c,
			params,
			routePathname,
		)

		resource, err := provider.Get().Handle(requestContext)
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

func buildFragment(fragment *views.DynamicFragmentView, c echo.Context) (string, error) {
	provider := resources.Provider.Find(fragment.GetResourceName())

	if provider.IsNil() {
		return "", errors.New("resource not found")
	}

	routePathname := c.Request().Header.Get("Hardwire-Dynamic-Fragment-Request")
	if routePathname == "" {
		return "", errors.New("bad request")
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
		return "", err
	}
	if resource == nil {
		return "", errors.New("resource not found")
	}

	html, err := fragment.Build(resource)
	if err != nil {
		return "", err
	}

	return html, nil
}
