package resourceprovider

import (
	"github.com/ncpa0/hardwire/utils"
	. "github.com/ncpa0cpl/convenient-structures"
)

type DynamicResource struct {
	name    string
	handler func(c *DynamicRequestContext) (interface{}, error)
}

func (resource *DynamicResource) Handle(c *DynamicRequestContext) (interface{}, error) {
	return resource.handler(c)
}

type DRProvider struct {
	resources *Map[string, *DynamicResource]
}

func (provider *DRProvider) GET(name string, handler func(c *DynamicRequestContext) (interface{}, error)) {
	provider.resources.Set(name, &DynamicResource{
		name:    name,
		handler: handler,
	})
}

func (provider *DRProvider) Find(name string) *utils.Option[DynamicResource] {
	res, _ := provider.resources.Get(name)
	return utils.NewOption(res)
}

var Provider *DRProvider = &DRProvider{
	resources: NewMap(map[string]*DynamicResource{}),
}
