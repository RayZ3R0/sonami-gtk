package ui

import (
	"codeberg.org/dergs/tidalwave/internal/router"
	"codeberg.org/dergs/tidalwave/internal/secrets"
	"codeberg.org/dergs/tidalwave/internal/signals"
	"codeberg.org/dergs/tidalwave/internal/ui/components"
	"codeberg.org/dergs/tidalwave/internal/ui/components/lyrics"
	"codeberg.org/dergs/tidalwave/internal/ui/components/player"
	"codeberg.org/dergs/tidalwave/internal/ui/components/queue"
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gio"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func (w *Window) buildSidebarHeader() *gtk.Widget {
	windowTitle := WindowTitle("Tidal Wave", "")()
	router.NavigationCompleted.On(func(entry router.HistoryEntry) bool {
		schwifty.OnMainThreadOncePure(func() {
			windowTitle.SetSubtitle(entry.PageTitle)
			w.SetTitle("Tidal Wave - " + entry.PageTitle)
		})
		return signals.Continue
	})

	mainMenu := gio.NewMenu()
	if secrets.HasRefreshToken() {
		mainMenu.Append("Sign Out", "win.sign-out")
	} else {
		mainMenu.Append("Sign In", "win.sign-in")
	}
	mainMenu.Append("Set as default page", "win.set-as-default")
	mainMenu.Append("Preferences", "app.preferences")
	mainMenu.Append("About", "app.about")
	mainMenu.Append("Quit", "app.quit")
	defer mainMenu.Unref()

	return HeaderBar().
		DecorationLayout("icon").
		TitleWidget(Widget(&windowTitle.Widget)).
		ShowBackButton(false).
		ShowEndTitleButtons(false).
		CenteringPolicy(adw.CenteringPolicyStrictValue).
		PackEnd(
			MenuButton().
				IconName("menu-symbolic").
				MenuModel(&mainMenu.MenuModel),
			components.NewRouteButton("search").Icon("loupe-symbolic"),
		).
		ToGTK()
}

func (w *Window) buildSidebar() schwifty.ViewStack {
	return ViewStack().
		AddTitledWithIcon(player.NewPlayer(), "player", "Player", "music-note-outline-symbolic").
		AddTitledWithIcon(lyrics.NewLyricsPanel(), "lyrics", "Lyrics", "chat-bubble-text-symbolic").
		AddTitledWithIcon(queue.NewQueue(), "queue", "Queue", "music-queue-symbolic")
}

func (w *Window) buildSidebarFooter(viewStack *adw.ViewStack) *gtk.Widget {
	viewSwitcher := adw.NewViewSwitcher()
	viewSwitcher.SetPolicy(adw.ViewSwitcherPolicyWideValue)
	viewSwitcher.SetStack(viewStack)
	return &viewSwitcher.Widget
}
