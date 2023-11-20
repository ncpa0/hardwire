package htmxframework

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/ncpa0/htmx-framework/views"
)

func createPageViewHandler(view *views.PageView) func(c echo.Context) error {
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
