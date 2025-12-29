package ui

import (
	"codeberg.org/dergs/tidalwave/internal/router"
	"codeberg.org/dergs/tidalwave/internal/signals"
	"codeberg.org/dergs/tidalwave/internal/ui/components"
	"codeberg.org/dergs/tidalwave/internal/ui/components/player"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gio"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func (w *Window) buildSidebarHeader() *gtk.Widget {
	windowTitle := WindowTitle("Tidal Wave", "")()
	router.NavigationComplete.On(func(entry router.HistoryEntry) bool {
		windowTitle.SetSubtitle(entry.PageTitle)
		w.SetTitle("Tidal Wave - " + entry.PageTitle)
		return signals.Continue
	})

	mainMenu := gio.NewMenu()
	mainMenu.Append("Preferences", "app.preferences")
	mainMenu.Append("About", "app.about")

	menuButton := gtk.NewMenuButton()
	menuButton.SetIconName("menu-symbolic")
	menuButton.SetMenuModel(&mainMenu.MenuModel)

	searchButton := components.NewRouteButton("search")
	searchButton.SetIcon("loupe-symbolic")

	return HeaderBar().
		DecorationLayout("icon").
		TitleWidget(Widget(&windowTitle.Widget)).
		ShowBackButton(false).
		ShowEndTitleButtons(false).
		CenteringPolicy(adw.CenteringPolicyStrictValue).
		PackEnd(menuButton, searchButton).
		ToGTK()
}

func (w *Window) buildSidebar() *adw.ViewStack {
	viewStack := adw.NewViewStack()
	viewStack.AddTitledWithIcon(player.NewPlayer().ToGTK(), "player", "Player", "music-note-outline-symbolic")
	// viewStack.AddTitledWithIcon(components.NewLyricsPanel(), "lyrics", "Lyrics", "chat-bubble-text-symbolic")
	// viewStack.AddTitledWithIcon(gtk.NewSpinner(), "queue", "Queue", "music-queue-symbolic")
	return viewStack
}

func (w *Window) buildSidebarFooter(viewStack *adw.ViewStack) *gtk.Widget {
	viewSwitcher := adw.NewViewSwitcher()
	viewSwitcher.SetPolicy(adw.ViewSwitcherPolicyWideValue)
	viewSwitcher.SetStack(viewStack)
	return &viewSwitcher.Widget
}
