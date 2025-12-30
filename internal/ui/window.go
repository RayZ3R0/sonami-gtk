package ui

import (
	"codeberg.org/dergs/tidalwave/internal/g"
	"codeberg.org/dergs/tidalwave/internal/notifications"
	"codeberg.org/dergs/tidalwave/internal/router"
	"codeberg.org/dergs/tidalwave/internal/signals"
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gio"
	"github.com/jwijenbergh/puregotk/v4/glib"
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

	window.installActions()
	window.SetContent(window.build())
	window.SetTitle("Tidal Wave")
	window.SetIconName("logo")
	window.SetDefaultSize(1280, 720)

	router.Navigate("home", nil)

	return window
}

func (w *Window) build() *gtk.Widget {
	layout := adw.NewOverlaySplitView()
	layout.SetSidebar(w.buildSidebarLayout())
	layout.SetContent(w.buildContentLayout())
	layout.SetSidebarWidthFraction(0.4)
	layout.SetMaxSidebarWidth(367)
	layout.SetMinSidebarWidth(367)

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
		toast := adw.NewToast(title)
		toast.SetTimeout(3)

		toastLayout.AddToast(toast)

		return signals.Continue
	})

	return &toastLayout.Widget
}

func (w *Window) buildContentLayout() *gtk.Widget {
	toolbarView := adw.NewToolbarView()
	toolbarView.AddTopBar(w.buildContentHeader())

	router.OnNavigate.On(func(path string) bool {
		schwifty.OnMainThreadOnce(func(u uintptr) {
			toolbarView.SetContent(loadingView.ToGTK())
		}, 0)
		return signals.Continue
	})

	router.NavigationComplete.On(func(entry router.HistoryEntry) bool {
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
