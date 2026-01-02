package syntax

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func Widget(w *gtk.Widget) schwifty.Widget {
	return func() *schwifty.WrappedWidget {
		return &schwifty.WrappedWidget{Widget: *w}
	}
}

// WARN: Do not manage reference counting for a schwifty-managed widget. If you are not in control of the widget's lifecycle, use Widget() instead.
func ManagedWidget(w *gtk.Widget) schwifty.Widget {
	return managed("ManagedWidget", func() *schwifty.WrappedWidget {
		return &schwifty.WrappedWidget{Widget: *w}
	})
}
