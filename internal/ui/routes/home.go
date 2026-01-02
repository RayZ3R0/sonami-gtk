package routes

import (
	"context"

	"codeberg.org/dergs/tidalwave/internal/router"
	"codeberg.org/dergs/tidalwave/internal/ui/components"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gtk"
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

	body := VStack().Spacing(25).VMargin(20)
	for _, item := range homeFeed.Items {
		body = body.Append(components.ForPageItem(item))
	}

	controller := gtk.NewEventControllerScroll(gtk.EventControllerScrollVerticalValue)
	controller.SetPropagationPhase(gtk.PhaseCaptureValue)

	// vadj := scroll.GetVadjustment()
	// controller.ConnectScroll(func(dx, dy float64) (ok bool) {
	// 	if controller.CurrentEventState()&gdk.KEY_Shift_L != 0 {
	// 		return false
	// 	}

	// 	if controller.Unit() == gdk.ScrollUnitWheel {
	// 		vadj.SetValue(vadj.Value() + dy*vadj.StepIncrement())
	// 	} else {
	// 		vadj.SetValue(vadj.Value() + dy)
	// 	}
	// 	return false
	// })

	return &router.Response{
		PageTitle: "Home",
		View: ScrolledWindow().
			AddController(&controller.EventController).
			Child(body).
			Policy(gtk.PolicyNeverValue, gtk.PolicyAutomaticValue),
	}
}
