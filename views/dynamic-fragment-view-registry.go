package views

import (
	"github.com/ncpa0/hardwire/utils"
	. "github.com/ncpa0cpl/convenient-structures"
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

func (vr *DynamicFragmentViewRegistry) GetFragment(routePathname string) *utils.Option[DynamicFragmentView] {
	viewsIterator := vr.views.Iterator()
	for !viewsIterator.Done() {
		view, _ := viewsIterator.Next()
		if view.MatchesRoute(routePathname) {
			return utils.NewOption(view)
		}
	}

	return utils.Empty[DynamicFragmentView]()
}

func (vr *DynamicFragmentViewRegistry) ForEach(cb func(view *DynamicFragmentView) error) error {
	viewsIterator := vr.views.Iterator()
	for !viewsIterator.Done() {
		view, _ := viewsIterator.Next()
		err := cb(view)
		if err != nil {
			return err
		}
	}
	return nil
}
