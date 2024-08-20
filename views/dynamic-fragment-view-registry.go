package views

import (
	"github.com/ncpa0/hardwire/utils"
	. "github.com/ncpa0cpl/ezs"
)

type DynamicFragmentViewRegistry struct {
	views *Array[*DynamicFragmentView]
}

func NewDynamicFragmentViewRegistry() *DynamicFragmentViewRegistry {
	return &DynamicFragmentViewRegistry{
		views: &Array[*DynamicFragmentView]{},
	}
}

func (vr *DynamicFragmentViewRegistry) Register(view *DynamicFragmentView) {
	vr.views.Push(view)
}

func (vr *DynamicFragmentViewRegistry) GetFragmentById(id string) *utils.Option[DynamicFragmentView] {
	for view := range vr.views.Iter() {
		if view.id == id {
			return utils.NewOption(view)
		}
	}

	return utils.Empty[DynamicFragmentView]()
}

func (vr *DynamicFragmentViewRegistry) GetFragment(routePathname string) *utils.Option[DynamicFragmentView] {
	for view := range vr.views.Iter() {
		if view.MatchesRoute(routePathname) {
			return utils.NewOption(view)
		}
	}

	return utils.Empty[DynamicFragmentView]()
}

func (vr *DynamicFragmentViewRegistry) ForEach(cb func(view *DynamicFragmentView) error) error {
	for view := range vr.views.Iter() {
		err := cb(view)
		if err != nil {
			return err
		}
	}
	return nil
}
