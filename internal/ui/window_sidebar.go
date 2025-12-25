package ui

import (
	"codeberg.org/dergs/tidalwave/internal/router"
	"codeberg.org/dergs/tidalwave/internal/ui/components"
	"codeberg.org/dergs/tidalwave/internal/ui/signals"
	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func (w *Window) buildSidebarHeader() gtk.Widgetter {
	windowTitle := adw.NewWindowTitle("Tidal Wave", "")
	router.NavigationComplete.On(func(response *router.Response) bool {
		windowTitle.SetSubtitle(response.PageTitle)
		w.SetTitle("Tidal Wave - " + response.PageTitle)
		return signals.Continue
	})

	headerbar := adw.NewHeaderBar()
	headerbar.SetDecorationLayout("icon")
	headerbar.SetTitleWidget(windowTitle)
	headerbar.SetShowBackButton(false)
	headerbar.SetShowEndTitleButtons(false)
	headerbar.SetCenteringPolicy(adw.CenteringPolicyStrict)

	mainMenu := gio.NewMenu()
	mainMenu.Append("Preferences", "app.preferences")
	mainMenu.Append("About", "app.about")

	menuButton := gtk.NewMenuButton()
	menuButton.SetIconName("menu-symbolic")
	menuButton.SetMenuModel(&mainMenu.MenuModel)

	headerbar.PackEnd(menuButton)

	searchButton := components.NewRouteButton("search")
	searchButton.SetIcon("loupe-symbolic")

	btn2 := gtk.NewButtonFromIconName("loupe-symbolic")
	btn2.ConnectClicked(func() {
		router.Navigate("search", nil)
	})
	headerbar.PackEnd(searchButton)
	return headerbar
}

func (w *Window) buildSidebar() *adw.ViewStack {
	viewStack := adw.NewViewStack()
	viewStack.AddTitledWithIcon(components.NewPlayer(), "player", "Player", "music-note-outline-symbolic")
	viewStack.AddTitledWithIcon(components.NewLyricsPanel(), "lyrics", "Lyrics", "chat-bubble-text-symbolic")
	viewStack.AddTitledWithIcon(gtk.NewSpinner(), "queue", "Queue", "music-queue-symbolic")
	return viewStack
}

func (w *Window) buildSidebarFooter(viewStack *adw.ViewStack) gtk.Widgetter {
	viewSwitcher := adw.NewViewSwitcher()
	viewSwitcher.SetPolicy(adw.ViewSwitcherPolicyWide)
	viewSwitcher.SetStack(viewStack)
	return viewSwitcher
}
