package hardwire

import (
	"net/http"

	echo "github.com/labstack/echo/v4"
	config "github.com/ncpa0/hardwire/configuration"
	"github.com/ncpa0/hardwire/utils"
	"github.com/ncpa0/hardwire/views"
)

type View interface {
	Render(c echo.Context) (*views.RenderedView, error)
}

func createResponse(c echo.Context, view View) error {
	ifNoneMatch := c.Request().Header.Get("If-None-Match")
	renderResult, err := view.Render(c)

	if err != nil {
		return utils.HandleError(c, err)
	}

	if renderResult.Etag != "" && ifNoneMatch == renderResult.Etag {
		return c.NoContent(http.StatusNotModified)
	}

	if renderResult.Etag != "" {
		c.Response().Header().Set("ETag", renderResult.Etag)
	}

	return c.HTML(http.StatusOK, renderResult.Html)
}

func createPageViewHandler(view *views.PageView, conf *config.Configuration) func(c echo.Context) error {
	if view.Metadata.ShouldRedirect {
		return func(c echo.Context) error {
			err := c.Redirect(http.StatusMovedPermanently, view.Metadata.RedirectURL)
			if err != nil {
				return err
			}
			if conf.BeforeResponse != nil {
				return conf.BeforeResponse(c)
			}
			return nil
		}
	}

	return func(c echo.Context) error {
		selector := c.Request().Header.Get("HX-Target")

		c.Response().Header().Set("Vary", "HX-Target")
		if !view.IsDynamic() {
			c.Response().Header().Set(
				"Cache-Control",
				config.GenerateCacheHeaderForStaticRoute(),
			)
		} else {
			c.Response().Header().Set(
				"Cache-Control",
				config.GenerateCacheHeaderForDynamicRoute(),
			)
		}

		if selector != "" {
			child := view.QuerySelector("#" + selector)

			if !child.IsNil() {
				err := createResponse(c, child.Get())
				if err != nil {
					return err
				}
				if conf.BeforeResponse != nil {
					return conf.BeforeResponse(c)
				}
				return nil
			}
		}

		err := createResponse(c, view)
		if err != nil {
			return err
		}
		if conf.BeforeResponse != nil {
			return conf.BeforeResponse(c)
		}
		return nil
	}
}
