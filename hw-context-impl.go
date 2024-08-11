package hardwire

import (
	"errors"
	"fmt"

	echo "github.com/labstack/echo/v4"
	hw "github.com/ncpa0/hardwire/hw-context"
	resources "github.com/ncpa0/hardwire/resources"
	"github.com/ncpa0/hardwire/utils"
)

type HwContext struct{}

func (ctx *HwContext) GetResourceHandler(e echo.Context, resourceKey string) (func(rootPath string, params map[string]string) (interface{}, error), error) {
	entry, found := resources.GetResource(resourceKey)
	if !found {
		e.String(404, "Invalid request")
		return nil, fmt.Errorf("Resource od key '%s' not found", resourceKey)
	}
	handler := resources.GetResourceHandler(entry)

	return func(rootPath string, params map[string]string) (interface{}, error) {
		dynReqCtx := resources.NewDynamicRequestContext(e, params, rootPath)
		return handler(dynReqCtx)
	}, nil
}

func (ctx *HwContext) BuildFragment(ectx echo.Context, fragment hw.BuildableFragment) (string, error) {
	resKey := fragment.ResourceKeys()[0]
	isValidResource := resources.HasResource(resKey)

	if !isValidResource {
		return "", errors.New("resource not found")
	}

	routePathname := ectx.Request().Header.Get("Hardwire-Dynamic-Fragment-Request")
	if routePathname == "" {
		return "", errors.New("bad request")
	}

	hxCurrentUrl := ectx.Request().Header.Get("Hx-Current-Url")
	params := utils.ParseUrlParams(routePathname, hxCurrentUrl)

	handler, err := HardwireContext.GetResourceHandler(ectx, resKey)
	if err != nil {
		return "", err
	}

	resource, err := handler(routePathname, params)
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
