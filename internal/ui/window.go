package ui

import (
	"codeberg.org/dergs/tidalwave/internal/g"
	"codeberg.org/dergs/tidalwave/internal/notifications"
	"codeberg.org/dergs/tidalwave/internal/router"
	"codeberg.org/dergs/tidalwave/internal/settings"
	"codeberg.org/dergs/tidalwave/internal/signals"
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gio"
	"github.com/jwijenbergh/puregotk/v4/glib"
	"github.com/jwijenbergh/puregotk/v4/gobject"
	"github.com/jwijenbergh/puregotk/v4/gtk"

	_ "codeberg.org/dergs/tidalwave/internal/ui/routes"
)

type Window struct {
	*adw.ApplicationWindow
}

var loadingView = syntax.Clamp().MaximumSize(50).Child(syntax.Spinner())

func NewWindow(app *adw.Application) *Window {
	window := &Window{
		ApplicationWindow: adw.NewApplicationWindow(&app.Application),
	}
	injector.Singleton(func() *adw.ApplicationWindow {
		return window.ApplicationWindow
	})
	injector.Singleton(func() *gtk.Window {
		return &window.Window
	})

	window.installActions()
	window.SetContent(window.build())
	window.SetTitle("Tidal Wave")
	window.SetIconName("logo")
	window.SetDefaultSize(settings.General().GetWindowWidth(), settings.General().GetWindowHeight())
	// For some reason the bindings do not allow to specify which property
	window.ConnectNotify(g.Ptr(func(gobject.Object, uintptr) {
		if window.GetHeight() > 0 {
			settings.General().SetWindowHeight(window.GetHeight())
		}
		if window.GetWidth() > 0 {
			settings.General().SetWindowWidth(window.GetWidth())
		}
	}))

	router.Navigate("home")

	return window
}

func (w *Window) build() *gtk.Widget {
	layout := adw.NewOverlaySplitView()
	layout.SetSidebar(w.buildSidebarLayout())
	layout.SetContent(w.buildContentLayout())
	layout.SetSidebarWidthFraction(0.4)
	layout.SetMaxSidebarWidth(420)
	layout.SetMinSidebarWidth(420)

	sidebarAction := gio.NewSimpleActionStateful("toggle-sidebar", nil, glib.NewVariantBoolean(true))
	sidebarAction.ConnectActivate(g.Ptr(func(action gio.SimpleAction, _ uintptr) {
		newState := !action.GetState().GetBoolean()
		action.SetState(glib.NewVariantBoolean(newState))
		layout.SetShowSidebar(newState)
	}))
	w.AddAction(sidebarAction)
	w.GetApplication().SetAccelsForAction("win.toggle-sidebar", []string{"<Ctrl>B"})

	toastLayout := adw.NewToastOverlay()
	toastLayout.SetChild(&layout.Widget)
	layout.Unref()

	notifications.OnToast.On(func(title string) bool {
		schwifty.OnMainThreadOncePure(func() {
			toast := adw.NewToast(title)
			toast.SetTimeout(3)

			toastLayout.AddToast(toast)
		})
		return signals.Continue
	})

	return &toastLayout.Widget
}

func (w *Window) buildContentLayout() *gtk.Widget {
	toolbarView := adw.NewToolbarView()
	toolbarView.AddTopBar(w.buildContentHeader())

	router.NavigationStarted.On(func(path string) bool {
		schwifty.OnMainThreadOnce(func(u uintptr) {
			toolbarView.SetContent(loadingView.ToGTK())
		}, 0)
		return signals.Continue
	})

	router.NavigationCompleted.On(func(entry router.HistoryEntry) bool {
		toolbarView.SetContent(entry.View)
		return signals.Continue
	})

	return &toolbarView.Widget
}

func (w *Window) buildSidebarLayout() *gtk.Widget {
	toolbarView := adw.NewToolbarView()
	toolbarView.AddTopBar(w.buildSidebarHeader())
	viewStack := w.buildSidebar()
	toolbarView.SetContent(&viewStack.Widget)

	box := gtk.NewCenterBox()
	box.SetCenterWidget(w.buildSidebarFooter(viewStack))
	box.SetMarginBottom(6)
	box.SetMarginStart(7)
	box.SetMarginEnd(7)
	box.SetMarginTop(6)

	toolbarView.AddBottomBar(&box.Widget)
	toolbarView.SetBottomBarStyle(adw.ToolbarFlatValue)
	return &toolbarView.Widget
}
