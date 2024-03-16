package hardwire

import (
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"

	echo "github.com/labstack/echo/v4"
	"github.com/ncpa0/hardwire/utils"
	"github.com/ncpa0/hardwire/views"
	. "github.com/ncpa0cpl/convenient-structures"
)

type ActionContext struct {
	Echo      echo.Context
	isHandled bool
	// list of islands that have been written
	// to the response so far
	updatedIslands []string
}

func (actx *ActionContext) Reload() {
	pageViewRegistry := views.GetPageViewRegistry()

	actx.isHandled = true
	currentUrl, err := url.Parse(actx.Echo.Request().Header.Get("HX-Current-URL"))
	if err != nil {
		actx.Echo.NoContent(http.StatusResetContent)
		return
	}

	view := pageViewRegistry.GetView(currentUrl.Path)
	if view.IsNil() {
		actx.Echo.NoContent(http.StatusResetContent)
		return
	}

	renderResult, err := view.Get().Render(actx.Echo)

	if err != nil {
		utils.HandleError(actx.Echo, err)
		return
	}

	actx.Echo.Response().Header().Set("HX-Retarget", "body")
	actx.Echo.HTML(200, renderResult.Html)
}

func (actx *ActionContext) Redirect(to string) {
	pageViewRegistry := views.GetPageViewRegistry()

	actx.isHandled = true
	if to[0] != '/' {
		to = "/" + to
	}

	view := pageViewRegistry.GetView(to)
	if view.IsNil() {
		actx.Echo.Redirect(http.StatusSeeOther, to)
		return
	}

	renderResult, err := view.Get().Render(actx.Echo)

	if err != nil {
		utils.HandleError(actx.Echo, err)
		return
	}

	actx.Echo.Response().Header().Set("HX-Push-Url", to)
	actx.Echo.Response().Header().Set("HX-Retarget", "body")
	actx.Echo.HTML(200, renderResult.Html)
}

func (actx *ActionContext) UpdateIslands(islandsIDs ...string) {
	allIslands := views.GetIslands()
	dynFragments := views.GetDynamicFragmentViewRegistry()
	for _, islandID := range islandsIDs {
		if slices.Contains(actx.updatedIslands, islandID) {
			continue
		}

		ok, island := allIslands.Find(func(island *views.Island, _ int) bool {
			return island.ID == islandID
		})

		if !ok {
			actx.Echo.Logger().Error("island not found: ", islandID)
			return
		}

		fragment := dynFragments.GetFragmentById(island.FragmentID)

		if fragment.IsNil() {
			actx.Echo.Logger().Error("fragment not found: ", island.FragmentID, ", required by island: ", island.ID)
			return
		}

		html, err := buildFragment(fragment.Get(), actx.Echo)
		if err != nil {
			actx.Echo.Logger().Error("error building fragment: ", err)
			return
		}

		actx.Echo.Response().Write([]byte(fmt.Sprintf("\n\n<div id=\"%s\" hx-swap-oob=\"true\">%s</div>", island.ID, html)))
		actx.updatedIslands = append(actx.updatedIslands, islandID)
		actx.isHandled = true
	}
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
	actx := &ActionContext{
		Echo: ctx,
	}
	err = action.handler(body, actx)
	if err != nil {
		return err
	}

	islandsHeader := ctx.Request().Header.Get("Hardwire-Islands-Update")
	allIslands := views.GetIslands()
	dynFragments := views.GetDynamicFragmentViewRegistry()
	if islandsHeader != "" && utils.IsStatusPositive(ctx.Response().Status) {
		islandsIDs := strings.Split(islandsHeader, ";")
		for _, islandID := range islandsIDs {
			if slices.Contains(actx.updatedIslands, islandID) {
				continue
			}

			ok, island := allIslands.Find(func(i1 *views.Island, i2 int) bool {
				return i1.ID == islandID
			})

			if !ok {
				ctx.Logger().Error("island not found: ", islandID)
				return echo.ErrNotFound
			}

			fragment := dynFragments.GetFragmentById(island.FragmentID)

			if fragment.IsNil() {
				ctx.Logger().Error("fragment not found: ", island.FragmentID, ", required by island: ", island.ID)
				return echo.ErrNotFound
			}

			html, err := buildFragment(fragment.Get(), ctx)
			if err != nil {
				ctx.Logger().Error("error building fragment: ", err)
				return echo.ErrInternalServerError
			}

			ctx.Response().Write([]byte(fmt.Sprintf("\n\n<div id=\"%s\" hx-swap-oob=\"true\">%s</div>", island.ID, html)))
			actx.updatedIslands = append(actx.updatedIslands, islandID)
			actx.isHandled = true
		}
	}

	if !actx.isHandled {
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
