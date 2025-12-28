package syntax

import (
	"log/slog"
	"runtime"

	"codeberg.org/dergs/tidalwave/internal/g"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var logger = slog.With("library", "schwifty")

type unrefable interface {
	ConnectDestroy(cb *func(gtk.Widget)) uint32
	GoPointer() uintptr
	Unref()
}

func finalize[T unrefable](tag string, widgetToBeFinalized T) {
	logger.Debug("now managing cleanup for object", "ptr", widgetToBeFinalized.GoPointer(), "tag", tag)
	widgetToBeFinalized.ConnectDestroy(g.Ptr(func(a gtk.Widget) {
		logger.Debug("object was destroyed by GTK", "ptr", a.GoPointer(), "tag", tag)
	}))
	runtime.SetFinalizer(widgetToBeFinalized, func(finalizedWidget T) {
		logger.Debug("releasing reference on schwifty-managed object, GTK may still hold a reference to this object", "tag", tag, "ptr", finalizedWidget.GoPointer())
		finalizedWidget.Unref()
	})
}

func managed[InnerType unrefable, T func() InnerType](tag string, callback T) T {
	return func() InnerType {
		widget := callback()
		finalize(tag, widget)
		return widget
	}
}
