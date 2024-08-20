package resourceprovider

import (
	"fmt"
	"net/http"
	"net/url"
	"slices"

	echo "github.com/labstack/echo/v4"
	hw "github.com/ncpa0/hardwire/hw-context"
	"github.com/ncpa0/hardwire/utils"
	"github.com/ncpa0/hardwire/views"
	"github.com/ncpa0cpl/ezs"
)

type ActionContext struct {
	HwContext          hw.HardwireContext
	Echo               echo.Context
	wasResponseWritten bool
	// list of islands that have been written
	// to the response so far
	updatedIslands []string
}

func (actx *ActionContext) Reload() {
	pageViewRegistry := views.GetPageViewRegistry()

	actx.wasResponseWritten = true
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

	renderResult, err := view.Get().Render(actx.HwContext, actx.Echo)

	if err != nil {
		utils.HandleError(actx.Echo, err)
		return
	}

	actx.Echo.Response().Header().Set("HX-Retarget", "body")
	actx.Echo.HTML(200, renderResult.Html)
}

func (actx *ActionContext) Redirect(to string) {
	pageViewRegistry := views.GetPageViewRegistry()

	actx.wasResponseWritten = true
	if to[0] != '/' {
		to = "/" + to
	}

	view := pageViewRegistry.GetView(to)
	if view.IsNil() {
		actx.Echo.Redirect(http.StatusSeeOther, to)
		return
	}

	renderResult, err := view.Get().Render(actx.HwContext, actx.Echo)

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

		requiredResources := fragment.Get().ResourceKeys()
		resources := ezs.NewMap(map[string]interface{}{})

		for _, resourceKey := range requiredResources {
			res, err := actx.HwContext.GetResource(actx.Echo, resourceKey)
			if err != nil {
				actx.Echo.Logger().Error("error getting resource: ", err)
				return
			}
			resources.Set(resourceKey, res)
		}

		html, err := actx.HwContext.BuildFragment(fragment.Get(), resources)
		if err != nil {
			actx.Echo.Logger().Error("error building fragment: ", err)
			return
		}

		morphSwap := actx.Echo.Request().Header.Get("Hardwire-Htmx-Morph") == "true"

		swap := utils.OobSwap{
			Selector: "#" + island.ID,
		}
		if morphSwap {
			swap.Extension = "morph"
		}

		utils.SetChunkedEnc(actx.Echo)
		actx.Echo.Response().Write([]byte(fmt.Sprintf("\n<div hx-swap-oob=\"%s\">%s</div>", swap.Build(), html)))
		actx.Echo.Response().Flush()
		actx.updatedIslands = append(actx.updatedIslands, islandID)
		actx.wasResponseWritten = true
	}
}
