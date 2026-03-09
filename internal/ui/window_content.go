package ui

import (
	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/internal/router"
	"github.com/RayZ3R0/sonami-gtk/internal/signals"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
)

func (w *Window) buildContentHeader() *gtk.Widget {
	homeButton := components.NewRouteButton("home", true)
	homeButton.Title(gettext.Get("Home"))
	homeButton.Icon("go-home-symbolic")
	homeButton.TooltipText(gettext.Get("Navigate to Home"))

	exploreButton := components.NewRouteButton("explore", true)
	exploreButton.Title(gettext.Get("Explore"))
	exploreButton.Icon("compass2-symbolic")
	exploreButton.TooltipText(gettext.Get("Navigate to Explore"))

	collectionButton := components.NewRouteButton("my-collection", true)
	collectionButton.Title(gettext.Get("My Collection"))
	collectionButton.Icon("heart-outline-thick-symbolic")
	collectionButton.TooltipText(gettext.Get("Navigate to My Collection"))

	defaultToolbar := HStack(
		homeButton,
		exploreButton,
		collectionButton,
	).Spacing(3)()

	// We never want to delete the default toolbar. NEVER.
	defaultToolbar.Ref()

	headerbar := HeaderBar().
		ShowStartTitleButtons(false).
		PackStart(
			Button().
				IconName("dock-left-symbolic").
				ActionName("win.toggle-sidebar").
				TooltipText(gettext.Get("Toggle Sidebar")),
			Button().
				IconName("left-symbolic").
				ActionName("win.navigate-back").
				Visible(false).
				TooltipText(gettext.Get("Navigate Back")).
				ConnectConstruct(func(b *gtk.Button) {
					router.HistoryUpdated.On(func(history *router.History) bool {
						schwifty.OnMainThreadOncePure(func() {
							b.SetVisible(len(history.Entries) > 0)
						})
						return signals.Continue
					})
				}),
		).
		TitleWidget(defaultToolbar)()

	router.Navigation.On(func(event *router.NavigationEvent) bool {
		schwifty.OnMainThreadOncePure(func() {
			if event.Completed && event.Result.Toolbar != nil {
				headerbar.SetTitleWidget(event.Result.Toolbar)
			} else {
				headerbar.SetTitleWidget(&defaultToolbar.Widget)
			}
		})
		return signals.Continue
	})

	return &headerbar.Widget
}
