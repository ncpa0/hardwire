package resourceprovider

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/antchfx/xmlquery"
	echo "github.com/labstack/echo/v4"
	hw "github.com/ncpa0/hardwire/hw-context"
	"github.com/ncpa0/hardwire/utils"
	"github.com/ncpa0/hardwire/views"
	. "github.com/ncpa0cpl/ezs"
	Promise "github.com/ncpa0cpl/go_promise"
)

type AtomicRespWriter struct {
	actionCtx *ActionContext
	mutex     *sync.Mutex
}

func (arw *AtomicRespWriter) SendIslandUpdate(islandID string, data []byte) {
	arw.mutex.Lock()
	defer arw.mutex.Unlock()

	utils.SetChunkedEnc(arw.actionCtx.Echo)
	resp := arw.actionCtx.Echo.Response()
	resp.Write(data)
	resp.Flush()
	arw.actionCtx.updatedIslands = append(arw.actionCtx.updatedIslands, islandID)
	arw.actionCtx.wasResponseWritten = true
}

type QueuedIsland struct {
	Island            *views.Island
	Fragment          *views.DynamicFragmentView
	RequiredResources *Array[string]
}

func (qi *QueuedIsland) CanRender(readyRes *Map[string, interface{}]) bool {
	return qi.RequiredResources.Every(func(resKey string, idx int) bool {
		return readyRes.Has(resKey)
	})
}

type Action struct {
	Name    string
	Method  string
	NewBody func() interface{}
	Handler func(body interface{}, ctx *ActionContext) error
}

func NewAction[T interface{}](
	name string,
	method string,
	handler func(body *T, ctx *ActionContext) error,
) *Action {
	action := &Action{
		Name:   name,
		Method: method,
		NewBody: func() interface{} {
			return new(T)
		},
		Handler: func(body interface{}, ctx *ActionContext) error {
			return handler(body.(*T), ctx)
		},
	}

	return action
}

func (action *Action) Perform(hwContext hw.HardwireContext, ctx echo.Context) error {
	body := action.NewBody()
	err := ctx.Bind(body)
	if err != nil {
		return echo.ErrBadRequest
	}
	err = bindFormParams(ctx, body)
	if err != nil {
		return echo.ErrInternalServerError
	}
	actx := &ActionContext{
		HwContext: hwContext,
		Echo:      ctx,
	}
	err = action.Handler(body, actx)
	if err != nil {
		return err
	}

	islandIDs := utils.ParseHeaderList(
		ctx.Request().Header.Get("Hardwire-Islands-Update"),
	)
	itemKeys := utils.ParseHeaderList(
		ctx.Request().Header.Get("Hardwire-Dynamic-List-Patch"),
	)
	morphSwap := ctx.Request().Header.Get("Hardwire-Htmx-Morph") == "true"

	if !utils.IsStatusPositive(ctx.Response().Status) {
		return nil
	}

	atomicWriter := &AtomicRespWriter{
		actionCtx: actx,
		mutex:     &sync.Mutex{},
	}

	sendFragmentUpdate := func(island *views.Island, html string) error {
		swap := utils.OobSwap{}
		if morphSwap {
			swap.Extension = "morph"
		}

		strReader := strings.NewReader(html)
		fragmentNode, err := xmlquery.Parse(strReader)
		if err != nil {
			ctx.Logger().Error("error parsing fragment output html: ", err)
			return echo.ErrInternalServerError
		}

		if itemKeys.Length() == 0 || island.Type != "list" {
			fragmentNode = xmlquery.FindOne(
				fragmentNode, "//div[@data-frag-url]",
			)

			swap.Selector = "#" + island.ID
			utils.XmlNodeSetAttribute(
				fragmentNode,
				"hx-swap-oob",
				swap.Build(),
			)
			nodeHtml := utils.XmlNodeToString(fragmentNode)

			atomicWriter.SendIslandUpdate(island.ID, []byte(nodeHtml))
		} else {
			items := NewArray([]string{})
			for itemKey := range itemKeys.Iter() {
				itemNode, err := xmlquery.Query(
					fragmentNode, fmt.Sprintf("//div[@data-item-key=\"%s\"]", itemKey),
				)
				if err == nil && itemNode != nil {
					swap.Selector = fmt.Sprintf(
						".island_%s .dynamic-list-element[data-item-key='%s']",
						island.ID, itemKey,
					)
					utils.XmlNodeSetAttribute(
						itemNode,
						"hx-swap-oob",
						swap.Build(),
					)
					items.Push(utils.XmlNodeToString(itemNode))
				} else {
					swap := utils.OobSwap{
						Mode: "delete",
						Selector: fmt.Sprintf(
							".island_%s .dynamic-list-element[data-item-key='%s']",
							island.ID, itemKey,
						),
					}
					items.Push(fmt.Sprintf(
						"<div hx-swap-oob=\"%s\"></div>",
						swap.Build(),
					))
				}
			}
			atomicWriter.SendIslandUpdate(island.ID, []byte("\n"+strings.Join(items.ToSlice(), "\n")))
		}

		return nil
	}

	allIslands := views.GetIslands()
	islandsToUpdate := allIslands.Filter(func(island *views.Island, i int) bool {
		return Contains(islandIDs, island.ID) && !Contains(NewArray(actx.updatedIslands), island.ID)
	})

	dynFragments := views.GetDynamicFragmentViewRegistry()
	qis := MapTo(islandsToUpdate, func(island *views.Island) *QueuedIsland {
		fragment := dynFragments.GetFragmentById(island.FragmentID).Get()
		var requiredResources *Array[string]
		if fragment != nil {
			requiredResources = NewArray(fragment.ResourceKeys())
		}
		return &QueuedIsland{
			Island:            island,
			Fragment:          fragment,
			RequiredResources: requiredResources,
		}
	}).Filter(func(qi *QueuedIsland, i int) bool {
		if qi.Fragment == nil {
			ctx.Logger().Error("fragment not found: ", qi.Island.FragmentID, ", required by island: ", qi.Island.ID)
			return false
		}
		return true
	})

	allRequiredResources := NewMap(map[string]bool{})
	for qi := range qis.Iter() {
		for resKey := range qi.RequiredResources.Iter() {
			allRequiredResources.Set(resKey, true)
		}
	}

	qisMutex := sync.Mutex{}
	readyResources := NewMap(map[string]interface{}{})

	renderPossibleIslands := func() *Array[error] {
		qisMutex.Lock()
		qisBatch := NewArray([]*QueuedIsland{})
		qis = qis.Filter(func(qi *QueuedIsland, idx int) bool {
			if qi.CanRender(readyResources) {
				qisBatch.Push(qi)
				return false
			}
			return true
		})
		qisMutex.Unlock()

		ops := MapTo(qisBatch, func(qi *QueuedIsland) *Promise.Promise[interface{}] {
			return Promise.New(func() (interface{}, error) {
				html, err := actx.HwContext.BuildFragment(qi.Fragment, readyResources)
				if err != nil {
					ctx.Logger().Errorf(
						"error building fragment for island (%s): %s",
						qi.Island.ID, err.Error(),
					)
					return nil, err
				}
				err = sendFragmentUpdate(qi.Island, html)
				return nil, err
			})
		})

		res := NewArray(*Promise.AwaitAll(ops.ToSlice()))
		failedRes := res.Filter(func(res Promise.AwaitAllResult[interface{}], i int) bool {
			return res.Err != nil
		})

		return MapTo(failedRes, func(res Promise.AwaitAllResult[interface{}]) error {
			return res.Err
		})
	}

	ops := MapTo(allRequiredResources.Keys(), func(resKey string) *Promise.Promise[interface{}] {
		return Promise.New(func() (interface{}, error) {
			res, err := actx.HwContext.GetResource(ctx, resKey)
			if err != nil {
				return nil, err
			}
			readyResources.Set(resKey, res)
			errs := renderPossibleIslands()
			if errs.Length() > 0 {
				return nil, errors.New("error occurred when rendering some islands")
			}
			return nil, nil
		})
	})

	results := NewArray(*Promise.AwaitAll(ops.ToSlice()))
	errors := MapTo(results.Filter(func(res Promise.AwaitAllResult[interface{}], i int) bool {
		return res.Err != nil
	}), func(res Promise.AwaitAllResult[interface{}]) string {
		return res.Err.Error()
	})

	allFailed := results.Every(func(res Promise.AwaitAllResult[interface{}], i int) bool {
		return res.Err != nil
	})

	if allFailed {
		ctx.Logger().Error(
			"error occured when rendering some islands or obtaining resources",
			Join(errors, ", "),
		)
		return ctx.String(http.StatusInternalServerError, "error occurred when rendering islands")
	}

	someFailed := errors.Length() > 0
	if someFailed {
		ctx.Logger().Error(
			"error occured when rendering some islands or obtaining resources",
			Join(errors, ", "),
		)
		return nil
	}

	if !actx.wasResponseWritten {
		return ctx.NoContent(http.StatusNoContent)
	}

	return nil
}
