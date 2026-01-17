package syntax

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty/callback"
	"codeberg.org/dergs/tonearm/pkg/schwifty/tracking"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

type manageableObject interface {
	tracking.Trackable
}

type manageableWidget interface {
	tracking.Trackable
	ConnectDestroy(cb *func(gtk.Widget)) uint32
	ConnectMap(cb *func(gtk.Widget)) uint32
	ConnectRealize(cb *func(gtk.Widget)) uint32
	ConnectUnmap(cb *func(gtk.Widget)) uint32
	ConnectUnrealize(cb *func(gtk.Widget)) uint32
}

func managedWidget[InnerType manageableWidget, T func() InnerType](tag string, fn T) T {
	return func() InnerType {
		widget := fn()
		widget.ConnectDestroy(&callback.DestroyCallback)
		widget.ConnectMap(&callback.MapCallback)
		widget.ConnectRealize(&callback.RealizedCallback)
		widget.ConnectUnmap(&callback.UnmapCallback)
		widget.ConnectUnrealize(&callback.UnrealizedCallback)
		tracking.SetFinalizer(tag, widget)
		return widget
	}
}

func managedObject[InnerType manageableObject, T func() InnerType](tag string, fn T) T {
	return func() InnerType {
		widget := fn()
		tracking.SetFinalizer(tag, widget)
		return widget
	}
}
