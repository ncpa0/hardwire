package htmxframework

import (
	"net/http"

	"github.com/labstack/echo"
	config "github.com/ncpa0/htmx-framework/configuration"
	"github.com/ncpa0/htmx-framework/utils"
	"github.com/ncpa0/htmx-framework/views"
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

func createPageViewHandler(view *views.PageView) func(c echo.Context) error {
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
				return createResponse(c, child.Get())
			}
		}

		return createResponse(c, view)
	}
}
