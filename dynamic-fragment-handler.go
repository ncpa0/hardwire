package htmxframework

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/ncpa0/htmx-framework/views"
)

func createDynamicFragmentHandler(view *views.DynamicFragmentView) func(c echo.Context) error {
	return func(c echo.Context) error {
		provider := DynamicResourceProvider.find(view.GetResourceName())
		resource, err := provider.handler(c)
		if err != nil {
			return err
		}
		if resource == nil {
			return c.NoContent(http.StatusNotFound)
		}

		html, err := view.Build(resource)
		if err != nil {
			return err
		}

		c.Response().Header().Set("Content-Type", "text/html")
		return c.String(http.StatusOK, html)
	}
}
