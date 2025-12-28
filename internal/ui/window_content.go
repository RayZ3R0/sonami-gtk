package ui

import (
	"codeberg.org/dergs/tidalwave/internal/router"
	"codeberg.org/dergs/tidalwave/internal/signals"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func (w *Window) buildContentHeader() *gtk.Widget {
	headerbar := adw.NewHeaderBar()
	headerbar.SetShowStartTitleButtons(false)

	sidebarButton := gtk.NewButtonFromIconName("sidebar-show-symbolic")
	sidebarButton.SetActionName("win.toggle-sidebar")
	headerbar.PackStart(&sidebarButton.Widget)

	backButton := gtk.NewButtonFromIconName("left-symbolic")
	backButton.SetActionName("win.navigate-back")
	backButton.SetVisible(false)
	headerbar.PackStart(&backButton.Widget)

	router.HistoryUpdated.On(func(history *router.History) bool {
		backButton.SetVisible(history.Length() > 1)
		return signals.Continue
	})

	// routeButton := components.NewRouteButton("home")
	// routeButton.SetTitle("Home")
	// routeButton.SetIcon("go-home-symbolic")

	// routeButton2 := components.NewRouteButton("explore")
	// routeButton2.SetTitle("Explore")
	// routeButton2.SetIcon("compass2-symbolic")

	// routeButton3 := components.NewRouteButton("collection")
	// routeButton3.SetTitle("Collection")
	// routeButton3.SetIcon("library-symbolic")

	// defaultTopBar := gtk.NewBox(gtk.OrientationHorizontalValue, 3)
	// defaultTopBar.Append(routeButton.Button)
	// defaultTopBar.Append(routeButton2.Button)
	// defaultTopBar.Append(routeButton3.Button)
	// headerbar.SetTitleWidget(&defaultTopBar.Widget)

	// router.OnNavigate.On(func(path string) bool {
	// 	headerbar.SetTitleWidget(&defaultTopBar.Widget)
	// 	return signals.Continue
	// })

	// router.NavigationComplete.On(func(response *router.Response) bool {
	// 	if response.Toolbar != nil {
	// 		headerbar.SetTitleWidget(response.Toolbar)
	// 	} else {
	// 		headerbar.SetTitleWidget(&defaultTopBar.Widget)
	// 	}
	// 	return signals.Continue
	// })

	return &headerbar.Widget
}
