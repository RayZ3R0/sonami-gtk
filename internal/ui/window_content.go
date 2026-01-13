package ui

import (
	"codeberg.org/dergs/tonearm/internal/router"
	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/internal/ui/components"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func (w *Window) buildContentHeader() *gtk.Widget {
	homeButton := components.NewRouteButton("home")
	homeButton.Title("Home")
	homeButton.Icon("go-home-symbolic")

	exploreButton := components.NewRouteButton("explore")
	exploreButton.Title("Explore")
	exploreButton.Icon("compass2-symbolic")

	collectionButton := components.NewRouteButton("my-collection")
	collectionButton.Title("Collection")
	collectionButton.Icon("library-symbolic")

	defaultToolbar := HStack(
		Widget(&homeButton.Widget),
		Widget(&exploreButton.Widget),
		Widget(&collectionButton.Widget),
	).Spacing(3)()

	// We never want to delete the default toolbar. NEVER.
	defaultToolbar.Ref()

	headerbar := HeaderBar().
		ShowStartTitleButtons(false).
		PackStart(
			Button().
				IconName("dock-left-symbolic").
				ActionName("win.toggle-sidebar"),
			Button().
				IconName("left-symbolic").
				ActionName("win.navigate-back").
				Visible(false).
				ConnectConstruct(func(b *gtk.Button) {
					router.HistoryUpdated.On(func(history *router.History) bool {
						schwifty.OnMainThreadOncePure(func() {
							b.SetVisible(len(history.Entries) > 1)
						})
						return signals.Continue
					})
				}),
		).
		TitleWidget(defaultToolbar)()

	router.NavigationStarted.On(func(path string) bool {
		schwifty.OnMainThreadOnce(func(u uintptr) {
			headerbar.SetTitleWidget(&defaultToolbar.Widget)
		}, 0)
		return signals.Continue
	})

	router.NavigationCompleted.On(func(entry router.HistoryEntry) bool {
		schwifty.OnMainThreadOncePure(func() {
			if entry.Toolbar != nil {
				headerbar.SetTitleWidget(entry.Toolbar)
			} else {
				headerbar.SetTitleWidget(&defaultToolbar.Widget)
			}
		})
		return signals.Continue
	})

	return &headerbar.Widget
}
