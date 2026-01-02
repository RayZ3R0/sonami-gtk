package syntax

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty/callback"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/tracking"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

type manageableWidget interface {
	tracking.Trackable
	ConnectRealize(cb *func(gtk.Widget)) uint32
	ConnectUnrealize(cb *func(gtk.Widget)) uint32
}

func managed[InnerType manageableWidget, T func() InnerType](tag string, fn T) T {
	return func() InnerType {
		widget := fn()
		widget.ConnectRealize(&callback.RealizedCallback)
		widget.ConnectUnrealize(&callback.UnrealizedCallback)
		tracking.SetFinalizer(tag, widget)
		return widget
	}
}
