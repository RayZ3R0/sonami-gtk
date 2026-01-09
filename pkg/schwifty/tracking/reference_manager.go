package tracking

import (
	"runtime"

	"codeberg.org/dergs/tidalwave/pkg/schwifty/callback"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var (
	shouldLogLifecycle = false
)

var onDestroy = func(a gtk.Widget) {
	if shouldLogLifecycle {
		logger.Debug("object was destroyed by GTK", "ptr", a.GoPointer())
	}
	callback.CallbackHandler[any](a.Object, "destroy", a)
	callback.DeleteCallbacks(a)
	Untrack(a.GoPointer())
}

type Trackable interface {
	ConnectDestroy(cb *func(gtk.Widget)) uint32
	GoPointer() uintptr
	Unref()
}

func SetFinalizer[T Trackable](tag string, widgetToBeFinalized T) {
	ptr := widgetToBeFinalized.GoPointer()
	Track(ptr, tag)
	if shouldLogLifecycle {
		logger.Debug("now managing cleanup for object", "ptr", ptr, "tag", tag)
	}
	widgetToBeFinalized.ConnectDestroy(&onDestroy)
	runtime.SetFinalizer(widgetToBeFinalized, func(finalizedWidget T) {
		if shouldLogLifecycle {
			logger.Debug("releasing reference on schwifty-managed object, GTK may still hold a reference to this object", "tag", tag, "ptr", finalizedWidget.GoPointer())
		}
		finalizedWidget.Unref()
		TrackGC(ptr)
	})
}

func SetLogLifecycle(enabled bool) {
	shouldLogLifecycle = enabled
}
