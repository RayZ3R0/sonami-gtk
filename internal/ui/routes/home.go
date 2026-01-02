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

func Home() *router.Response {
	tidal := injector.MustInject[*tidalapi.TidalAPI]()
	homeFeed, err := tidal.V2.Home.Feed.Static(context.Background())
	if err != nil {
		return router.FromError("Home", err)
	}

	body := VStack().Spacing(25).VMargin(20)
	for _, item := range homeFeed.Items {
		body = body.Append(components.ForPageItem(item))
	}

	return &router.Response{
		PageTitle: "Home",
		View: ScrolledWindow().
			Child(body).
			Policy(gtk.PolicyNeverValue, gtk.PolicyAutomaticValue),
	}
}
