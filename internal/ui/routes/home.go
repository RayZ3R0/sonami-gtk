package routes

import (
	"context"

	"codeberg.org/dergs/tidalwave/internal/router"
	"codeberg.org/dergs/tidalwave/internal/ui/components/dynamic"
	. "codeberg.org/dergs/tidalwave/pkg/gui"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/infinytum/injector"
)

func init() {
	router.Register("home", Home)
}

func Home(params router.Params) *router.Response {
	tidal := injector.MustInject[*tidalapi.TidalAPI]()
	homeFeed, err := tidal.V2.Home.Feed.Static(context.Background())
	if err != nil {
		return router.FromError("Home", err)
	}

	children := make([]gtk.Widgetter, 0)
	for _, item := range homeFeed.Items {
		children = append(children, dynamic.ForPageItem(item))
	}

	scroll := gtk.NewScrolledWindow()
	scroll.SetPolicy(gtk.PolicyNever, gtk.PolicyAutomatic)
	scroll.SetChild(VStack(children...).Spacing(25).VMargin(20))

	controller := gtk.NewEventControllerScroll(gtk.EventControllerScrollVertical)
	controller.SetPropagationPhase(gtk.PhaseCapture)
	scroll.AddController(controller)

	vadj := scroll.VAdjustment()
	controller.ConnectScroll(func(dx, dy float64) (ok bool) {
		if controller.CurrentEventState()&gdk.KEY_Shift_L != 0 {
			return false
		}

		if controller.Unit() == gdk.ScrollUnitWheel {
			vadj.SetValue(vadj.Value() + dy*vadj.StepIncrement())
		} else {
			vadj.SetValue(vadj.Value() + dy)
		}
		return false
	})

	return &router.Response{
		PageTitle: "Home",
		View:      scroll,
	}
}
