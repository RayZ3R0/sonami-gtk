package schwifty

import (
	gtkbindings "codeberg.org/dergs/tonearm/pkg/schwifty/bindings/gtk"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

//go:generate go run codeberg.org/dergs/tonearm/pkg/schwifty/gen Widget *WrappedWidget
type WrappedWidget struct {
	gtk.Widget
}

type BaseWidgetable = gtkbindings.BaseWidgetable

type Widgetable[T any] = gtkbindings.Widgetable[T]
