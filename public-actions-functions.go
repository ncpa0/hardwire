package hardwire

import (
	resources "github.com/ncpa0/hardwire/resources"
)

func RegisterPostAction[T interface{}](
	resource *resources.ResourceEntry,
	name string,
	action func(body *T, ctx *resources.ActionContext) error,
) {
	resources.RegisterPostAction(resource, name, action)
}

func RegisterPutAction[T interface{}](
	resource *resources.ResourceEntry,
	name string,
	action func(body *T, ctx *resources.ActionContext) error,
) {
	resources.RegisterPutAction(resource, name, action)
}

func RegisterPatchAction[T interface{}](
	resource *resources.ResourceEntry,
	name string,
	action func(body *T, ctx *resources.ActionContext) error,
) {
	resources.RegisterPatchAction(resource, name, action)
}

func RegisterDeleteAction[T interface{}](
	resource *resources.ResourceEntry,
	name string,
	action func(body *T, ctx *resources.ActionContext) error,
) {
	resources.RegisterDeleteAction(resource, name, action)
}
