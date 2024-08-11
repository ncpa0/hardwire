package utils

import echo "github.com/labstack/echo/v4"

func ParamMap(c echo.Context) map[string]string {
	params := make(map[string]string)

	for _, param := range c.ParamNames() {
		params[param] = c.Param(param)
	}

	return params
}
