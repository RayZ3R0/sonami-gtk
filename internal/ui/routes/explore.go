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
	router.Register("explore", ExploreMain)
	router.Register("explore/:page", Explore)
}

func ExploreMain() *router.Response {
	return Explore("explore")
}

func Explore(pageName string) *router.Response {
	tidal := injector.MustInject[*tidalapi.TidalAPI]()
	page, err := tidal.V1.Pages.Page(context.Background(), pageName)
	if err != nil {
		return router.FromError("Explore", err)
	}

	body := VStack().Spacing(25).VMargin(20)
	for _, row := range page.Rows {
		for _, module := range row.Modules {
			body = body.Append(VStack(
				components.ForModule(module),
			))
		}
	}

	return &router.Response{
		PageTitle: page.Title,
		View: ScrolledWindow().
			Child(body).
			Policy(gtk.PolicyNeverValue, gtk.PolicyAutomaticValue),
	}
}
