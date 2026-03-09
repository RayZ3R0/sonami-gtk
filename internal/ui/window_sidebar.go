package ui

import (
	"codeberg.org/puregotk/puregotk/v4/adw"
	"codeberg.org/puregotk/puregotk/v4/gio"
	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/internal/router"
	"github.com/RayZ3R0/sonami-gtk/internal/signals"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components/sidebar/lyrics"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components/sidebar/player"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components/sidebar/queue"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
)

func (w *Window) buildSidebarHeader() *gtk.Widget {
	windowTitle := WindowTitle("Sonami", "")()
	router.Navigation.On(func(entry *router.NavigationEvent) bool {
		if entry.Completed {
			schwifty.OnMainThreadOncePure(func() {
				w.SetTitle("Sonami - " + entry.Result.PageTitle)
			})
		}
		return signals.Continue
	})

	mainMenu := gio.NewMenu()
	mainMenu.Append(gettext.Get("Set as Default Page"), "win.set-as-default")
	mainMenu.Append(gettext.Get("Keyboard Shortcuts"), "app.shortcuts")
	mainMenu.Append(gettext.Get("Preferences"), "app.preferences")
	mainMenu.Append(gettext.Get("About Sonami"), "app.about")
	mainMenu.Append(gettext.Get("Quit"), "app.quit")

	return HeaderBar().
		TitleWidget(Widget(&windowTitle.Widget)).
		ShowBackButton(false).
		ShowEndTitleButtons(false).
		CenteringPolicy(adw.CenteringPolicyStrictValue).
		PackStart(
			components.NewRouteButton("search", false).Icon("loupe-symbolic").TooltipText(gettext.Get("Search")),
		).
		PackEnd(
			MenuButton().
				IconName("menu-symbolic").
				MenuModel(&mainMenu.MenuModel).
				TooltipText(gettext.Get("Main Menu")).ConnectConstruct(func(mb *gtk.MenuButton) {
				menuAction := gio.NewSimpleAction("main-menu", nil)
				menuAction.ConnectActivate(new(func(action gio.SimpleAction, parameter uintptr) {
					mb.Popup()
				}))
				w.AddAction(menuAction)
				w.GetApplication().SetAccelsForAction("win.main-menu", []string{"F10"})
			}),
		).
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
