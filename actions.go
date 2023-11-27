package htmxframework

import (
	"net/http"
	"net/url"

	"github.com/labstack/echo"
	"github.com/ncpa0/htmx-framework/utils"
	"github.com/ncpa0/htmx-framework/views"
	. "github.com/ncpa0cpl/convenient-structures"
)

type ActionContext struct {
	Echo      echo.Context
	isHandled bool
}

func (ctx *ActionContext) Reload() {
	pageViewRegistry := views.GetPageViewRegistry()

	ctx.isHandled = true
	currentUrl, err := url.Parse(ctx.Echo.Request().Header.Get("HX-Current-URL"))
	if err != nil {
		ctx.Echo.NoContent(http.StatusResetContent)
		return
	}

	view := pageViewRegistry.GetView(currentUrl.Path)
	if view.IsNil() {
		ctx.Echo.NoContent(http.StatusResetContent)
		return
	}

	renderResult, err := view.Get().Render(ctx.Echo)

	if err != nil {
		utils.HandleError(ctx.Echo, err)
		return
	}

	ctx.Echo.Response().Header().Set("HX-Retarget", "body")
	ctx.Echo.HTML(200, renderResult.Html)
}

func (ctx *ActionContext) Redirect(to string) {
	pageViewRegistry := views.GetPageViewRegistry()

	ctx.isHandled = true
	if to[0] != '/' {
		to = "/" + to
	}

	view := pageViewRegistry.GetView(to)
	if view.IsNil() {
		ctx.Echo.Redirect(http.StatusSeeOther, to)
		return
	}

	renderResult, err := view.Get().Render(ctx.Echo)

	if err != nil {
		utils.HandleError(ctx.Echo, err)
		return
	}

	ctx.Echo.Response().Header().Set("HX-Push-Url", to)
	ctx.Echo.Response().Header().Set("HX-Retarget", "body")
	ctx.Echo.HTML(200, renderResult.Html)
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

var actions *Array[iaction] = &Array[iaction]{}

func PostAction[Body interface{}](
	actionName string,
	handler func(body *Body, ctx *ActionContext) error,
) {
	action := &action[Body]{
		method:  "POST",
		name:    actionName,
		handler: handler,
	}

	actions.Push(action)
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

	actions.Push(action)
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

	actions.Push(action)
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

	actions.Push(action)
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

	actions.Push(action)
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
	actIterator := actions.Iterator()
	for !actIterator.Done() {
		action, _ := actIterator.Next()
		server.Add(
			action.GetMethod(),
			"/__actions/"+action.GetName(),
			func(ctx echo.Context) error {
				return action.Perform(ctx)
			},
		)
	}
}
