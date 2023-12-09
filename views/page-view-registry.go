package views

import (
	"github.com/ncpa0/hardwire/utils"
	. "github.com/ncpa0cpl/convenient-structures"
)

type PageViewRegistry struct {
	views *Array[*PageView]
}

func NewViewRegistry() *PageViewRegistry {
	return &PageViewRegistry{
		views: &Array[*PageView]{},
	}
}

func (vr *PageViewRegistry) Register(view *PageView) {
	vr.views.Push(view)
}

func (vr *PageViewRegistry) GetView(routePathname string) *utils.Option[PageView] {
	viewIterator := vr.views.Iterator()
	for !viewIterator.Done() {
		view, _ := viewIterator.Next()
		if view.MatchesRoute(routePathname) {
			return utils.NewOption(view)
		}
	}

	return utils.Empty[PageView]()
}

func (vr *PageViewRegistry) ForEach(cb func(view *PageView)) {
	viewIterator := vr.views.Iterator()
	for !viewIterator.Done() {
		view, _ := viewIterator.Next()
		cb(view)
	}
}
