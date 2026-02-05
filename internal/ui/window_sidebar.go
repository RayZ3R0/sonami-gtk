package ui

import (
	"strings"

	"codeberg.org/dergs/tonearm/internal/g"
	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/router"
	"codeberg.org/dergs/tonearm/internal/secrets"
	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/internal/ui/components"
	"codeberg.org/dergs/tonearm/internal/ui/components/lyrics"
	"codeberg.org/dergs/tonearm/internal/ui/components/player"
	"codeberg.org/dergs/tonearm/internal/ui/components/queue"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gio"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var decorationLayoutState = state.NewStateful("icon,appmenu:close")

func (w *Window) buildSidebarHeader() *gtk.Widget {
	windowTitle := WindowTitle("Tonearm", "")()
	router.NavigationCompleted.On(func(entry router.HistoryEntry) bool {
		schwifty.OnMainThreadOncePure(func() {
			windowTitle.SetSubtitle(entry.PageTitle)
			w.SetTitle("Tonearm - " + entry.PageTitle)
		})
		return signals.Continue
	})

	mainMenu := gio.NewMenu()
	mainMenu.Append(gettext.Get("Sign In"), "win.sign-in")
	mainMenu.Append(gettext.Get("Set as default page"), "win.set-as-default")
	mainMenu.Append(gettext.Get("Keyboard Shortcuts"), "app.shortcuts")
	mainMenu.Append(gettext.Get("Preferences"), "app.preferences")
	mainMenu.Append(gettext.Get("About Tonearm"), "app.about")
	mainMenu.Append(gettext.Get("Quit"), "app.quit")

	secrets.SignedInChanged.On(func(signedIn bool) bool {
		mainMenu.Remove(0)
		if signedIn {
			mainMenu.Insert(0, gettext.Get("Sign Out"), "win.sign-out")
		} else {
			mainMenu.Insert(0, gettext.Get("Sign In"), "win.sign-in")
		}
		return signals.Continue
	})

	settings := gtk.SettingsGetDefault()
	settings.ConnectSignal("notify::gtk-decoration-layout", g.Ptr(func() {
		updateDecorationLayout()
	}))
	updateDecorationLayout()

	return HeaderBar().
		BindDecorationLayout(decorationLayoutState).
		TitleWidget(Widget(&windowTitle.Widget)).
		ShowBackButton(false).
		ShowEndTitleButtons(false).
		CenteringPolicy(adw.CenteringPolicyStrictValue).
		PackEnd(
			MenuButton().
				IconName("menu-symbolic").
				MenuModel(&mainMenu.MenuModel).
				TooltipText(gettext.Get("Main Menu")).ConnectConstruct(func(mb *gtk.MenuButton) {
				menuAction := gio.NewSimpleAction("main-menu", nil)
				menuAction.ConnectActivate(g.Ptr(func(action gio.SimpleAction, parameter uintptr) {
					mb.Popup()
				}))
				w.AddAction(menuAction)
				w.GetApplication().SetAccelsForAction("win.main-menu", []string{"F10"})
			}),
			components.NewRouteButton("search").Icon("loupe-symbolic").TooltipText(gettext.Get("Search")),
			components.NewRouteButton("feed").Icon("bell-outline-symbolic").TooltipText(gettext.Get("Feed")),
		).
		ConnectDestroy(func(w gtk.Widget) {
			settings.Unref()
		}).
		ToGTK()
}

func (w *Window) buildSidebar() schwifty.ViewStack {
	return ViewStack().
		AddTitledWithIcon(player.NewPlayer(), "player", gettext.Get("Player"), "music-note-outline-symbolic").
		AddTitledWithIcon(lyrics.NewLyricsPanel(), "lyrics", gettext.Get("Lyrics"), "chat-bubble-text-symbolic").
		AddTitledWithIcon(queue.NewQueue(), "queue", gettext.Get("Queue"), "music-queue-symbolic")
}

func (w *Window) buildSidebarFooter(viewStack *adw.ViewStack) *gtk.Widget {
	viewSwitcher := adw.NewViewSwitcher()
	viewSwitcher.SetPolicy(adw.ViewSwitcherPolicyWideValue)
	viewSwitcher.SetStack(viewStack)
	return &viewSwitcher.Widget
}

func updateDecorationLayout() {
	settings := gtk.SettingsGetDefault()
	defer settings.Unref()

	configured := settings.GetPropertyGtkDecorationLayout()
	splits := strings.Split(configured, ":")
	left := splits[0]
	right := ""
	if len(splits) > 1 {
		right = splits[1]
	}

	if left == "appmenu" {
		decorationLayoutState.SetValue("icon," + left + ":" + right)
	} else if right == "appmenu" {
		decorationLayoutState.SetValue(left + ":" + right + ",icon")
	} else {
		decorationLayoutState.SetValue(configured)
	}
}
