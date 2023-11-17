package htmxframework

import (
	"net/http"
	"net/url"

	"github.com/labstack/echo"
)

type ActionContext struct {
	Echo      echo.Context
	isHandled bool
}

func (ctx *ActionContext) Reload() {
	ctx.isHandled = true
	currentUrl, err := url.Parse(ctx.Echo.Request().Header.Get("HX-Current-URL"))
	if err != nil {
		ctx.Echo.NoContent(http.StatusResetContent)
		return
	}

	view := viewRegistry.GetView(currentUrl.Path)
	if view.IsNil() {
		ctx.Echo.NoContent(http.StatusResetContent)
		return
	}

	ctx.Echo.Response().Header().Set("HX-Retarget", "body")
	ctx.Echo.String(200, view.Get().GetNode().ToHtml())
}

func (ctx *ActionContext) Redirect(to string) {
	ctx.isHandled = true
	if to[0] != '/' {
		to = "/" + to
	}

	view := viewRegistry.GetView(to)
	if view.IsNil() {
		ctx.Echo.Redirect(http.StatusSeeOther, to)
		return
	}

	ctx.Echo.Response().Header().Set("HX-Push-Url", to)
	ctx.Echo.Response().Header().Set("HX-Retarget", "body")
	ctx.Echo.String(200, view.Get().GetNode().ToHtml())
}

func (ctx *ActionContext) ReloadFragment(fragmentName string) {
	panic("not implemented")
	// ctx.isHandled = true
	// currentUrl, err := url.Parse(ctx.Echo.Request().Header.Get("HX-Current-URL"))
	// if err != nil {
	// 	ctx.Echo.NoContent(http.StatusResetContent)
	// }

	// view := viewRegistry.GetView(currentUrl.Path)
	// if view.IsNil() {
	// 	ctx.Echo.NoContent(http.StatusResetContent)
	// }

	// ctx.Echo.Response().Header().Set("HX-Retarget", "body")
	// ctx.Echo.String(200, view.Get().GetNode().ToHtml())
}

var actions []iaction = make([]iaction, 0)

func PostAction[Body interface{}](
	actionName string,
	handler func(body *Body, ctx *ActionContext) error,
) {
	action := &action[Body]{
		method:  "POST",
		name:    actionName,
		handler: handler,
	}

	actions = append(actions, action)
}

func PutAction[Body interface{}](
	actionName string,
	handler func(body *Body, ctx *ActionContext) error,
) {
	action := &action[Body]{
		method:  "PUT",
		name:    actionName,
		handler: handler,
	}

	actions = append(actions, action)
}

func PatchAction[Body interface{}](
	actionName string,
	handler func(body *Body, ctx *ActionContext) error,
) {
	action := &action[Body]{
		method:  "PATCH",
		name:    actionName,
		handler: handler,
	}

	actions = append(actions, action)
}

func DeleteAction[Body interface{}](
	actionName string,
	handler func(body *Body, ctx *ActionContext) error,
) {
	action := &action[Body]{
		method:  "DELETE",
		name:    actionName,
		handler: handler,
	}

	actions = append(actions, action)
}

func GetAction[Body interface{}](
	actionName string,
	handler func(body *Body, ctx *ActionContext) error,
) {
	action := &action[Body]{
		method:  "GET",
		name:    actionName,
		handler: handler,
	}

	actions = append(actions, action)
}

type action[Body interface{}] struct {
	method  string
	name    string
	handler func(body *Body, ctx *ActionContext) error
}

type iaction interface {
	GetMethod() string
	GetName() string
	Perform(ctx echo.Context) error
}

func (action *action[Body]) GetMethod() string {
	return action.method
}

func (action *action[Body]) GetName() string {
	return action.name
}

func (action *action[Body]) Perform(ctx echo.Context) error {
	body := new(Body)
	err := ctx.Bind(body)
	if err != nil {
		return echo.ErrBadRequest
	}
	actionContext := &ActionContext{
		Echo: ctx,
	}
	err = action.handler(body, actionContext)
	if err != nil {
		return err
	}

	if !actionContext.isHandled {
		return ctx.NoContent(http.StatusNoContent)
	}

	return nil
}

func registerActionHandlers(server *echo.Echo) {
	for _, action := range actions {

		server.Add(
			action.GetMethod(),
			"/__actions/"+action.GetName(),
			func(ctx echo.Context) error {
				return action.Perform(ctx)
			},
		)
	}
}
