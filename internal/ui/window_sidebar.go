package ui

import (
	"codeberg.org/dergs/tonearm/internal/ui/sidebar"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func (w *Window) buildSidebarLayout() *gtk.Widget {
	toolbarView := adw.NewToolbarView()
	toolbarView.AddTopBar(sidebar.BuildSidebarHeader())
	viewStack := sidebar.BuildSidebar()()
	toolbarView.SetContent(&viewStack.Widget)

	toolbarView.AddBottomBar(CenterBox().
		CenterWidget(sidebar.BuildSidebarFooter(viewStack)).
		HMargin(7).VMargin(6).ToGTK())
	toolbarView.SetBottomBarStyle(adw.ToolbarFlatValue)
	return &toolbarView.Widget
}
