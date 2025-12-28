package syntax

import (
	"log/slog"
	"runtime"

	"codeberg.org/dergs/tidalwave/pkg/schwifty/tracking"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var logger = slog.With("library", "schwifty")

var onDestroy = func(a gtk.Widget) {
	logger.Debug("object was destroyed by GTK", "ptr", a.GoPointer())
	tracking.Untrack(a.GoPointer())
}

type unrefable interface {
	ConnectDestroy(cb *func(gtk.Widget)) uint32
	GoPointer() uintptr
	Unref()
}

func finalize[T unrefable](tag string, widgetToBeFinalized T) {
	ptr := widgetToBeFinalized.GoPointer()
	tracking.Track(ptr, tag)
	logger.Debug("now managing cleanup for object", "ptr", ptr, "tag", tag)
	widgetToBeFinalized.ConnectDestroy(&onDestroy)
	runtime.SetFinalizer(widgetToBeFinalized, func(finalizedWidget T) {
		logger.Debug("releasing reference on schwifty-managed object, GTK may still hold a reference to this object", "tag", tag, "ptr", finalizedWidget.GoPointer())
		finalizedWidget.Unref()
		tracking.TrackGC(ptr)
	})
}

func managed[InnerType unrefable, T func() InnerType](tag string, callback T) T {
	return func() InnerType {
		widget := callback()
		finalize(tag, widget)
		return widget
	}
}
