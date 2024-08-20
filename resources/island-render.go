package resourceprovider

import (
	"fmt"
	"strings"
	"sync"

	"github.com/antchfx/xmlquery"
	echo "github.com/labstack/echo/v4"
	"github.com/ncpa0/hardwire/utils"
	"github.com/ncpa0/hardwire/views"
	. "github.com/ncpa0cpl/ezs"
	Promise "github.com/ncpa0cpl/go_promise"
)

type RespWriter interface {
	Write(islandID string, data []byte) error
}

func sendIslandUpdate(
	ctx echo.Context, writer RespWriter,
	island *views.Island, html string,
	morphSwap bool, itemKeys *Array[string],
) error {
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

		writer.Write(island.ID, []byte(nodeHtml))
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
		writer.Write(island.ID, []byte("\n"+strings.Join(items.ToSlice(), "\n")))
	}

	return nil
}

func renderIslands(
	actx *ActionContext,
	mutex *sync.Mutex,
	islands *Array[*QueuedIsland],
	resources *Map[string, interface{}],
	writer RespWriter,
	morphSwap bool,
	itemKeys *Array[string],
) *Array[error] {
	mutex.Lock()
	qisBatch := NewArray([]*QueuedIsland{})
	islands = islands.Filter(func(qi *QueuedIsland, idx int) bool {
		if qi.CanRender(resources) {
			qisBatch.Push(qi)
			return false
		}
		return true
	})
	mutex.Unlock()

	ops := MapTo(qisBatch, func(qi *QueuedIsland) *Promise.Promise[interface{}] {
		return Promise.New(func() (interface{}, error) {
			html, err := actx.HwContext.BuildFragment(qi.Fragment, resources)
			if err != nil {
				actx.Echo.Logger().Errorf(
					"error building fragment for island (%s): %s",
					qi.Island.ID, err.Error(),
				)
				return nil, err
			}
			err = sendIslandUpdate(
				actx.Echo, writer, qi.Island, html, morphSwap, itemKeys,
			)
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
