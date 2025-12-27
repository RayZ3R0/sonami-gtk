package ui

import (
	"context"

	"codeberg.org/dergs/tidalwave/internal/router"
	_ "codeberg.org/dergs/tidalwave/internal/ui/routes"
	"codeberg.org/dergs/tidalwave/internal/ui/signals"
	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotkit/app"
)

type Window struct {
	*adw.ApplicationWindow
}

func NewWindow(ctx context.Context) *Window {
	appInstance := app.FromContext(ctx)
	window := &Window{
		ApplicationWindow: adw.NewApplicationWindow(appInstance.Application),
	}

	window.SetContent(window.build())
	window.SetTitle("Tidal Wave")
	window.SetIconName("logo")
	window.SetDefaultSize(1280, 720)

	router.Navigate("home", nil)

	return window
}

func (w *Window) build() gtk.Widgetter {
	layout := adw.NewOverlaySplitView()
	layout.SetSidebar(w.buildSidebarLayout())
	layout.SetContent(w.buildContentLayout())
	layout.SetSidebarWidthFraction(0.4)
	layout.SetMaxSidebarWidth(367)
	layout.SetMinSidebarWidth(367)

	sidebarAction := gio.NewSimpleActionStateful("toggle-sidebar", nil, glib.NewVariantBoolean(true))
	sidebarAction.ConnectActivate(func(parameter *glib.Variant) {
		sidebarAction.SetState(glib.NewVariantBoolean(!sidebarAction.State().Boolean()))
		layout.SetShowSidebar(sidebarAction.State().Boolean())
	})
	w.AddAction(sidebarAction)
	w.Application().SetAccelsForAction("win.toggle-sidebar", []string{"<Ctrl>B"})

	navigateBackAction := gio.NewSimpleAction("navigate-back", nil)
	navigateBackAction.ConnectActivate(func(parameter *glib.Variant) {
		router.Back()
	})
	w.AddAction(navigateBackAction)
	w.Application().SetAccelsForAction("win.navigate-back", []string{"<Alt>Left"})

	toastLayout := adw.NewToastOverlay()
	toastLayout.SetChild(layout)

	signals.OnDisplayToast.On(func(val string) bool {
		toast := adw.NewToast(val)
		toast.SetTimeout(2)

		toastLayout.AddToast(toast)

		return signals.Continue
	})

	return toastLayout
}

func (w *Window) buildContentLayout() gtk.Widgetter {
	toolbarView := adw.NewToolbarView()
	toolbarView.AddTopBar(w.buildContentHeader())

	router.OnNavigate.On(func(path string) bool {
		spinner := gtk.NewSpinner()
		spinner.SetSpinning(true)
		spinner.Start()

		clamp := adw.NewClamp()
		clamp.SetMaximumSize(50)
		clamp.SetChild(spinner)
		toolbarView.SetContent(clamp)
		return signals.Continue
	})

	router.NavigationComplete.On(func(response *router.Response) bool {
		toolbarView.SetContent(response.View)
		return signals.Continue
	})

	return toolbarView
}

func (w *Window) buildSidebarLayout() gtk.Widgetter {
	toolbarView := adw.NewToolbarView()
	toolbarView.AddTopBar(w.buildSidebarHeader())
	viewStack := w.buildSidebar()
	toolbarView.SetContent(viewStack)

	box := gtk.NewCenterBox()
	box.SetCenterWidget(w.buildSidebarFooter(viewStack))
	box.SetMarginBottom(6)
	box.SetMarginStart(7)
	box.SetMarginEnd(7)
	box.SetMarginTop(6)

	toolbarView.AddBottomBar(box)
	toolbarView.SetBottomBarStyle(adw.ToolbarFlat)
	return toolbarView
}
