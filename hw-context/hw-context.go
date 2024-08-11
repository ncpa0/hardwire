package hwcontext

import (
	"github.com/labstack/echo/v4"
)

type BuildableFragment interface {
	ResourceKeys() []string
	Build(resources ...interface{}) (string, error)
}

type HardwireContext interface {
	GetResourceHandler(echoCtx echo.Context, resourceKey string) (
		func(rootPath string, params map[string]string) (interface{}, error),
		error,
	)
	BuildFragment(echoCtx echo.Context, fragment BuildableFragment) (string, error)
}
