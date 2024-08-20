package hwcontext

import (
	"github.com/labstack/echo/v4"
	. "github.com/ncpa0cpl/ezs"
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
	GetResource(echoCtx echo.Context, resourceKey string) (interface{}, error)
	BuildFragment(fragment BuildableFragment, resources *Map[string, interface{}]) (string, error)
}
