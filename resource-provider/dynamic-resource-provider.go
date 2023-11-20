package resourceprovider

import (
	"github.com/ncpa0/htmx-framework/utils"
)

type DynamicResource struct {
	name    string
	handler func(c *DynamicRequestContext) (interface{}, error)
}

func (resource *DynamicResource) Handle(c *DynamicRequestContext) (interface{}, error) {
	return resource.handler(c)
}

type DRProvider struct {
	resources map[string]*DynamicResource
}

func (provider *DRProvider) GET(name string, handler func(c *DynamicRequestContext) (interface{}, error)) {
	provider.resources[name] = &DynamicResource{
		name:    name,
		handler: handler,
	}
}

func (provider *DRProvider) Find(name string) *utils.Option[DynamicResource] {
	res := provider.resources[name]
	return utils.NewOption(res)
}

var Provider *DRProvider = &DRProvider{
	resources: make(map[string]*DynamicResource),
}
