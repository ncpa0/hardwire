package views

import "github.com/ncpa0/htmx-framework/utils"

type DynamicFragmentViewRegistry struct {
	views []*DynamicFragmentView
}

func NewDynamicFragmentViewRegistry() *DynamicFragmentViewRegistry {
	return &DynamicFragmentViewRegistry{
		views: []*DynamicFragmentView{},
	}
}

func (vr *DynamicFragmentViewRegistry) Register(view *DynamicFragmentView) {
	vr.views = append(vr.views, view)
}

func (vr *DynamicFragmentViewRegistry) GetFragment(routePathname string) *utils.Option[DynamicFragmentView] {
	for _, view := range vr.views {
		if view.MatchesRoute(routePathname) {
			return utils.NewOption(view)
		}
	}

	return utils.Empty[DynamicFragmentView]()
}

func (vr *DynamicFragmentViewRegistry) ForEach(cb func(view *DynamicFragmentView)) {
	for _, view := range vr.views {
		cb(view)
	}
}
