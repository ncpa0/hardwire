package resourceprovider

import (
	"fmt"

	echo "github.com/labstack/echo/v4"
	"github.com/ncpa0/hardwire/configuration"
	hw "github.com/ncpa0/hardwire/hw-context"
)

func RegisterPostAction[T interface{}](
	resource *ResourceEntry,
	name string,
	action func(body *T, ctx *ActionContext) error,
) {
	resource.actions.Push(NewAction(name, "POST", action))
}

func RegisterPutAction[T interface{}](
	resource *ResourceEntry,
	name string,
	action func(body *T, ctx *ActionContext) error,
) {
	resource.actions.Push(NewAction(name, "PUT", action))
}

func RegisterPatchAction[T interface{}](
	resource *ResourceEntry,
	name string,
	action func(body *T, ctx *ActionContext) error,
) {
	resource.actions.Push(NewAction(name, "PATCH", action))
}

func RegisterDeleteAction[T interface{}](
	resource *ResourceEntry,
	name string,
	action func(body *T, ctx *ActionContext) error,
) {
	resource.actions.Push(NewAction(name, "DELETE", action))
}

func MountActionEndpoints(hwContext hw.HardwireContext, server *echo.Echo) {
	ResourceReg.resources.ForEach(func(resourceKey string, entry *ResourceEntry) {
		entry.actions.ForEach(func(action *Action, idx int) {
			endpointPath := fmt.Sprintf("/__resources/%s/actions/%s", resourceKey, action.Name)
			if configuration.Current.DebugMode {
				fmt.Printf(
					"Adding action endpoint: %s\n",
					endpointPath,
				)
			}
			server.Add(
				action.Method,
				endpointPath,
				func(ctx echo.Context) error {
					return action.Perform(hwContext, ctx)
				},
			)
		})
	})
}
