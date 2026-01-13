package schwifty

import "github.com/jwijenbergh/puregotk/v4/gtk"

//go:generate go run codeberg.org/dergs/tonearm/pkg/schwifty/gen Widget *WrappedWidget
type WrappedWidget struct {
	gtk.Widget
}

type BaseWidgetable interface {
	ToGTK() *gtk.Widget
}

type Widgetable[T any] interface {
	BaseWidgetable
	AddController(controller *gtk.EventController) T
	CSS(css string) T
	Focusable(focusable bool) T
	FocusOnClick(focusOnClick bool) T
	HAlign(align gtk.Align) T
	HExpand(expand bool) T
	HMargin(horizontal int) T
	Margin(margin int) T
	MarginBottom(bottom int) T
	MarginEnd(end int) T
	MarginStart(start int) T
	MarginTop(top int) T
	Opacity(opacity float64) T
	Overflow(overflow gtk.Overflow) T
	VAlign(align gtk.Align) T
	VExpand(expand bool) T
	Visible(visible bool) T
	VMargin(vertical int) T
}
