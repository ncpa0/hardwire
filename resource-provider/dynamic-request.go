package resourceprovider

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/ncpa0/htmx-framework/utils"
)

type DynamicRequestContext struct {
	Echo          echo.Context
	params        map[string]string
	routePathname string
}

func NewDynamicRequestContext(echo echo.Context, params map[string]string, routePathname string) *DynamicRequestContext {
	return &DynamicRequestContext{
		Echo:          echo,
		params:        params,
		routePathname: routePathname,
	}
}

// Returns the route's URL parameter value
func (ctx *DynamicRequestContext) GetParam(key string) string {
	return ctx.params[key]
}

// Returns the route's URL parameter value
func (ctx *DynamicRequestContext) GetRoutePath() string {
	return ctx.routePathname
}

func (ctx *DynamicRequestContext) Err(code int, message string) *ResourceRequestError {
	return &ResourceRequestError{
		errType: "error",
		RequestError: utils.RequestError{
			Code: code,
			Data: message,
		},
	}
}

func (ctx *DynamicRequestContext) Redirect(to string) *ResourceRequestError {
	return &ResourceRequestError{
		errType: "redirect",
		RequestError: utils.RequestError{
			Code: http.StatusSeeOther,
			Data: to,
		},
	}
}

type ResourceRequestError struct {
	utils.RequestError
	errType string
}

func (err *ResourceRequestError) SendResponse(c echo.Context) error {
	switch err.errType {
	case "error":
		return c.String(err.Code, err.Data)
	case "redirect":
		if c.Request().Header.Get("Hx-Request") != "" {
			c.Response().Header().Set("Hx-Redirect", err.Data)
			return c.NoContent(200)
		}
		return c.Redirect(err.Code, err.Data)
	}

	return c.NoContent(err.Code)
}
