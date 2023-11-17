package views

import "github.com/ncpa0/htmx-framework/utils"

type ViewRegistry struct {
	views []*View
}

func NewViewRegistry() *ViewRegistry {
	return &ViewRegistry{
		views: []*View{},
	}
}

func (vr *ViewRegistry) Register(view *View) {
	vr.views = append(vr.views, view)
}

func (vr *ViewRegistry) GetView(routePathname string) *utils.Option[View] {
	for _, view := range vr.views {
		if view.MatchesRoute(routePathname) {
			return utils.NewOption(view)
		}
	}

	return utils.Empty[View]()
}

func (vr *ViewRegistry) ForEach(cb func(view *View)) {
	for _, view := range vr.views {
		cb(view)
	}
}
