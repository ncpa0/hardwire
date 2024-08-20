package resourceprovider

import (
	"github.com/ncpa0/hardwire/utils"
	. "github.com/ncpa0cpl/ezs"
)

type ResourceEntry struct {
	name     string
	resource Resource[interface{}]
	actions  Array[*Action]
}

func (entry *ResourceEntry) findAction(method string, name string) (bool, *Action) {
	return entry.actions.Find(func(action *Action, _ int) bool {
		return action.Method == method && action.Name == name
	})
}

type ResourceRegistry struct {
	resources *Map[string, *ResourceEntry]
}

func (reg *ResourceRegistry) Register(name string, resource Resource[interface{}]) *ResourceEntry {
	entry := &ResourceEntry{
		name:     name,
		resource: resource,
		actions:  Array[*Action]{},
	}
	reg.resources.Set(name, entry)
	return entry
}

func (reg *ResourceRegistry) find(name string) (*ResourceEntry, bool) {
	return reg.resources.Find(func(key string, value *ResourceEntry) bool {
		return key == name
	})
}

func (reg *ResourceRegistry) findByResource(resource Resource[interface{}]) *utils.Option[ResourceEntry] {
	entry, found := reg.resources.Find(func(key string, value *ResourceEntry) bool {
		return value.resource == resource
	})
	if found {
		return utils.NewOption(entry)
	} else {
		return utils.Empty[ResourceEntry]()
	}
}

var ResourceReg *ResourceRegistry = &ResourceRegistry{
	resources: NewMap(map[string]*ResourceEntry{}),
}

func HasResource(name string) bool {
	return ResourceReg.resources.Has(name)
}

func GetResource(name string) (*ResourceEntry, bool) {
	return ResourceReg.resources.Find(func(key string, value *ResourceEntry) bool {
		return key == name
	})
}

func GetResourceHandler(entry *ResourceEntry) func(c *DynamicRequestContext) (interface{}, error) {
	return entry.resource.Get
}
