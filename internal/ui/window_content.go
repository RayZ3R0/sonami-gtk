package ui

import (
	"codeberg.org/dergs/tidalwave/internal/router"
	"codeberg.org/dergs/tidalwave/internal/signals"
	"codeberg.org/dergs/tidalwave/internal/ui/components"
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func (w *Window) buildContentHeader() *gtk.Widget {

	sidebarButton := gtk.NewButtonFromIconName("sidebar-show-symbolic")
	sidebarButton.SetActionName("win.toggle-sidebar")

	backButton := gtk.NewButtonFromIconName("left-symbolic")
	backButton.SetActionName("win.navigate-back")
	backButton.SetVisible(false)
	router.HistoryUpdated.On(func(history *router.History) bool {
		backButton.SetVisible(len(history.Entries) > 1)
		return signals.Continue
	})

	homeButton := components.NewRouteButton("home")
	homeButton.SetTitle("Home")
	homeButton.SetIcon("go-home-symbolic")

	exploreButton := components.NewRouteButton("explore")
	exploreButton.SetTitle("Explore")
	exploreButton.SetIcon("compass2-symbolic")

	collectionButton := components.NewRouteButton("my-collection")
	collectionButton.SetTitle("Collection")
	collectionButton.SetIcon("library-symbolic")

	defaultToolbar := HStack(
		Widget(&homeButton.Widget),
		Widget(&exploreButton.Widget),
		Widget(&collectionButton.Widget),
	).Spacing(3)()

	// We never want to delete the default toolbar. NEVER.
	defaultToolbar.Ref()

	headerbar := HeaderBar().
		ShowStartTitleButtons(false).
		PackStart(sidebarButton, backButton).
		TitleWidget(defaultToolbar)()

	router.NavigationStarted.On(func(path string) bool {
		schwifty.OnMainThreadOnce(func(u uintptr) {
			headerbar.SetTitleWidget(nil)
		}, 0)
		return signals.Continue
	})

	router.NavigationCompleted.On(func(entry router.HistoryEntry) bool {
		if entry.Toolbar != nil {
			headerbar.SetTitleWidget(entry.Toolbar)
		} else {
			headerbar.SetTitleWidget(&defaultToolbar.Widget)
		}
		return signals.Continue
	})

	return &headerbar.Widget
}
