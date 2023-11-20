package views

import "github.com/ncpa0/htmx-framework/utils"

type PageViewRegistry struct {
	views []*PageView
}

func NewViewRegistry() *PageViewRegistry {
	return &PageViewRegistry{
		views: []*PageView{},
	}
}

func (vr *PageViewRegistry) Register(view *PageView) {
	vr.views = append(vr.views, view)
}

func (vr *PageViewRegistry) GetView(routePathname string) *utils.Option[PageView] {
	for _, view := range vr.views {
		if view.MatchesRoute(routePathname) {
			return utils.NewOption(view)
		}
	}

	return utils.Empty[PageView]()
}

func (vr *PageViewRegistry) ForEach(cb func(view *PageView)) {
	for _, view := range vr.views {
		cb(view)
	}
}
