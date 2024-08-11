package resourceprovider

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/antchfx/xmlquery"
	echo "github.com/labstack/echo/v4"
	hw "github.com/ncpa0/hardwire/hw-context"
	"github.com/ncpa0/hardwire/utils"
	"github.com/ncpa0/hardwire/views"
	. "github.com/ncpa0cpl/convenient-structures"
)

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

	swap := utils.OobSwap{}
	if morphSwap {
		swap.Extension = "morph"
	}

	allIslands := views.GetIslands()
	dynFragments := views.GetDynamicFragmentViewRegistry()
	if islandIDs.Length() > 0 && utils.IsStatusPositive(ctx.Response().Status) {
		for _, islandID := range islandIDs.ToSlice() {
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

			html, err := actx.HwContext.BuildFragment(ctx, fragment.Get())
			if err != nil {
				ctx.Logger().Error("error building fragment: ", err)
				return echo.ErrInternalServerError
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

				utils.SetChunkedEnc(ctx)
				ctx.Response().Write([]byte(nodeHtml))
				ctx.Response().Flush()
				actx.updatedIslands = append(actx.updatedIslands, islandID)
				actx.wasResponseWritten = true
			} else {
				items := Array[string]{}
				for _, itemKey := range itemKeys.ToSlice() {
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

				utils.SetChunkedEnc(ctx)
				ctx.Response().Write([]byte("\n" + strings.Join(items.ToSlice(), "\n")))
				ctx.Response().Flush()
				actx.updatedIslands = append(actx.updatedIslands, islandID)
				actx.wasResponseWritten = true
			}
		}
	}

	if !actx.wasResponseWritten {
		return ctx.NoContent(http.StatusNoContent)
	}

	return nil
}
