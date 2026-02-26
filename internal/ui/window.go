package ui

import (
	"codeberg.org/dergs/tonearm/internal/g"
	"codeberg.org/dergs/tonearm/internal/notifications"
	"codeberg.org/dergs/tonearm/internal/router"
	"codeberg.org/dergs/tonearm/internal/settings"
	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gio"
	"github.com/jwijenbergh/puregotk/v4/glib"
	"github.com/jwijenbergh/puregotk/v4/gobject"
	"github.com/jwijenbergh/puregotk/v4/gtk"

	_ "codeberg.org/dergs/tonearm/internal/ui/routes"
)

type Window struct {
	*adw.ApplicationWindow
}

var loadingView = g.Lazy(func() *gtk.Widget {
	widget := Clamp().MaximumSize(50).Child(Spinner()).ToGTK()
	widget.Ref()
	return widget
})

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
	window.installMouseClickHandler()

	window.SetContent(window.build())
	window.SetTitle("Tonearm")
	window.SetIconName("logo-symbolic")
	window.SetDefaultSize(settings.General().GetWindowWidth(), settings.General().GetWindowHeight())
	// For some reason the bindings do not allow to specify which property
	window.ConnectNotify(new(func(gobject.Object, uintptr) {
		if window.GetHeight() > 0 {
			settings.General().SetWindowHeight(window.GetHeight())
		}
		if window.GetWidth() > 0 {
			settings.General().SetWindowWidth(window.GetWidth())
		}
	}))

	router.Navigate(settings.General().DefaultPage())

	if !isStable() {
		window.AddCssClass("devel")
	}

	return window
}

func (w *Window) build() *gtk.Widget {
	layout := adw.NewOverlaySplitView()
	layout.SetSidebar(w.buildSidebarLayout())
	layout.SetContent(w.buildContentLayout())
	layout.SetSidebarWidthFraction(0.4)
	layout.SetMaxSidebarWidth(420)
	layout.SetMinSidebarWidth(320)

	sidebarAction := gio.NewSimpleActionStateful("toggle-sidebar", nil, glib.NewVariantBoolean(true))
	sidebarAction.ConnectActivate(new(func(action gio.SimpleAction, _ uintptr) {
		newState := !action.GetState().GetBoolean()
		action.SetState(glib.NewVariantBoolean(newState))
		layout.SetShowSidebar(newState)
	}))
	w.AddAction(sidebarAction)
	w.GetApplication().SetAccelsForAction("win.toggle-sidebar", []string{"<Ctrl>B", "F9"})

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

	router.Navigation.On(func(event *router.NavigationEvent) bool {
		schwifty.OnMainThreadOnce(func(u uintptr) {
			if event.Completed {
				toolbarView.SetContent(event.Result.View)
			} else {
				toolbarView.SetContent(loadingView())
			}
		}, 0)
		return signals.Continue
	})

	return &toolbarView.Widget
}
