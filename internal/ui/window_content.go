package ui

import (
	"codeberg.org/dergs/tidalwave/internal/router"
	"codeberg.org/dergs/tidalwave/internal/ui/components"
	"codeberg.org/dergs/tidalwave/internal/ui/signals"
	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func (w *Window) buildContentHeader() gtk.Widgetter {
	headerbar := adw.NewHeaderBar()
	headerbar.SetShowStartTitleButtons(false)

	sidebarButton := gtk.NewButtonFromIconName("sidebar-show-symbolic")
	sidebarButton.SetActionName("win.toggle-sidebar")
	headerbar.PackStart(sidebarButton)

	backButton := gtk.NewButtonFromIconName("left-symbolic")
	backButton.SetActionName("win.navigate-back")
	backButton.SetVisible(false)
	headerbar.PackStart(backButton)

	router.HistoryUpdated.On(func(history *router.History) bool {
		backButton.SetVisible(history.Length() > 1)
		return signals.Continue
	})

	routeButton := components.NewRouteButton("home")
	routeButton.SetTitle("Home")
	routeButton.SetIcon("go-home-symbolic")

	routeButton2 := components.NewRouteButton("explore")
	routeButton2.SetTitle("Explore")
	routeButton2.SetIcon("compass2-symbolic")

	routeButton3 := components.NewRouteButton("collection")
	routeButton3.SetTitle("Collection")
	routeButton3.SetIcon("library-symbolic")

	defaultTopBar := gtk.NewBox(gtk.OrientationHorizontal, 3)
	defaultTopBar.Append(routeButton)
	defaultTopBar.Append(routeButton2)
	defaultTopBar.Append(routeButton3)
	headerbar.SetTitleWidget(defaultTopBar)

	router.NavigationComplete.On(func(response *router.Response) bool {
		if response.Toolbar != nil {
			headerbar.SetTitleWidget(response.Toolbar)
		} else {
			headerbar.SetTitleWidget(defaultTopBar)
		}
		return signals.Continue
	})

	return headerbar
}
