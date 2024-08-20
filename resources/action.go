package resourceprovider

import (
	"errors"
	"net/http"
	"sync"

	echo "github.com/labstack/echo/v4"
	hw "github.com/ncpa0/hardwire/hw-context"
	"github.com/ncpa0/hardwire/utils"
	"github.com/ncpa0/hardwire/views"
	. "github.com/ncpa0cpl/ezs"
)

type AtomicRespWriter struct {
	actionCtx *ActionContext
	mutex     *sync.Mutex
}

func (arw *AtomicRespWriter) Write(islandID string, data []byte) error {
	arw.mutex.Lock()
	defer arw.mutex.Unlock()

	utils.SetChunkedEnc(arw.actionCtx.Echo)
	resp := arw.actionCtx.Echo.Response()
	_, err := resp.Write(data)
	if err != nil {
		return err
	}
	resp.Flush()
	arw.actionCtx.updatedIslands = append(arw.actionCtx.updatedIslands, islandID)
	arw.actionCtx.wasResponseWritten = true
	return nil
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

	allIslands := views.GetIslands()
	islandsToUpdate := allIslands.Filter(func(island *views.Island, i int) bool {
		return Contains(islandIDs, island.ID) && !Contains(NewArray(actx.updatedIslands), island.ID)
	})

	if islandsToUpdate.Length() == 0 {
		if !actx.wasResponseWritten {
			return ctx.NoContent(http.StatusNoContent)
		}

		return nil
	}

	dynFragments := views.GetDynamicFragmentViewRegistry()
	queuedIslands := MapTo(islandsToUpdate, func(island *views.Island) *QueuedIsland {
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

	// if there's only one island to update, do it in the current thread,
	// don't create new goroutines
	if queuedIslands.Length() == 1 {
		qIsland := queuedIslands.At(0)
		resources := NewMap(map[string]interface{}{})
		resKeys := NewArray(qIsland.Fragment.ResourceKeys())

		if resKeys.Length() == 1 {
			res, err := actx.HwContext.GetResource(ctx, resKeys.At(0))
			if err != nil {
				return err
			}
			resources.Set(resKeys.At(0), res)
		} else {
			// if there's more than one resource required, get them in parallel
			_, errs := utils.InParallel(
				resKeys.ToSlice(),
				func(key string) (interface{}, error) {
					res, err := actx.HwContext.GetResource(ctx, key)
					if err != nil {
						return nil, err
					}
					resources.Set(key, res)
					return nil, nil
				},
			)
			if len(errs) > 0 {
				errMsgs := MapTo(NewArray(errs), func(err error) string {
					return err.Error()
				})
				ctx.Logger().Error(
					"error occured when obtaining resources",
					Join(errMsgs, ", "),
				)
				return ctx.String(http.StatusInternalServerError, "error occurred when rendering")
			}
		}

		html, err := actx.HwContext.BuildFragment(qIsland.Fragment, resources)
		if err != nil {
			return err
		}
		err = sendIslandUpdate(
			actx.Echo, atomicWriter, qIsland.Island, html, morphSwap, itemKeys,
		)
		return err
	}

	allRequiredResourcesMap := NewMap(map[string]bool{})
	for qi := range queuedIslands.Iter() {
		for resKey := range qi.RequiredResources.Iter() {
			allRequiredResourcesMap.Set(resKey, true)
		}
	}
	allRequiredResources := allRequiredResourcesMap.Keys()

	respMutex := &sync.Mutex{}
	readyResources := NewMap(map[string]interface{}{})

	// if there's only one resource to retrieve, get it on the current thread
	// and render islands in separete goroutines
	if allRequiredResources.Length() == 1 {
		resKey := allRequiredResources.At(0)
		resource, err := actx.HwContext.GetResource(ctx, resKey)
		if err != nil {
			return err
		}
		readyResources.Set(resKey, resource)
		errs := renderIslands(
			actx, respMutex, queuedIslands, readyResources,
			atomicWriter, morphSwap, itemKeys,
		)
		if errs.Length() > 0 {
			return errors.New("error occurred when rendering some islands")
		}
		return nil
	}

	// retrive each resource in parallel, for each of them, once it's ready:
	// renderIslands will pick up all isalnds that can be rendered with resources
	// currently present, render them, write to the response and flush
	//
	// renderIslands will process each island in parallel
	_, errs := utils.InParallel(
		allRequiredResources.ToSlice(),
		func(resKey string) (interface{}, error) {
			res, err := actx.HwContext.GetResource(ctx, resKey)
			if err != nil {
				return nil, err
			}
			readyResources.Set(resKey, res)
			errs := renderIslands(
				actx, respMutex, queuedIslands, readyResources,
				atomicWriter, morphSwap, itemKeys,
			)
			if errs.Length() > 0 {
				return nil, errors.New("error occurred when rendering some islands")
			}
			return nil, nil
		},
	)

	errMsgs := MapTo(NewArray(errs), func(err error) string {
		return err.Error()
	})

	allFailed := len(errs) == allRequiredResources.Length()
	if allFailed {
		ctx.Logger().Error(
			"error occured when rendering some islands or obtaining resources",
			Join(errMsgs, ", "),
		)
		return ctx.String(http.StatusInternalServerError, "error occurred when rendering")
	}

	someFailed := errMsgs.Length() > 0
	if someFailed {
		ctx.Logger().Error(
			"error occured when rendering some islands or obtaining resources",
			Join(errMsgs, ", "),
		)
		return nil
	}

	if !actx.wasResponseWritten {
		return ctx.NoContent(http.StatusNoContent)
	}

	return nil
}
