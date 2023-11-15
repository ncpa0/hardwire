package htmxframework

import "github.com/labstack/echo"

type DynamicResource struct {
	name    string
	handler func(c echo.Context) (interface{}, error)
}

type DRProvider struct {
	resources map[string]*DynamicResource
}

func (provider *DRProvider) GET(name string, handler func(c echo.Context) (interface{}, error)) {
	provider.resources[name] = &DynamicResource{
		name:    name,
		handler: handler,
	}
}

func (provider *DRProvider) find(name string) *DynamicResource {
	return provider.resources[name]
}
