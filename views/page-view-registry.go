package views

import (
	"github.com/ncpa0/hardwire/utils"
	. "github.com/ncpa0cpl/ezs"
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
	for view := range vr.views.Iter() {
		if view.MatchesRoute(routePathname) {
			return utils.NewOption(view)
		}
	}

	return utils.Empty[PageView]()
}

func (vr *PageViewRegistry) ForEach(cb func(view *PageView) error) error {
	for view := range vr.views.Iter() {
		err := cb(view)
		if err != nil {
			return err
		}
	}
	return nil
}
