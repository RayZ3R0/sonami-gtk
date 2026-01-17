package tracking

import (
	"runtime"

	"codeberg.org/dergs/tonearm/pkg/schwifty/callback"
	"github.com/jwijenbergh/puregotk/v4/gobject"
)

var (
	shouldLogLifecycle = false
)

var onUnref = gobject.WeakNotify(func(userData uintptr, objPtr uintptr) {
	if shouldLogLifecycle {
		logger.Debug("object was destroyed by GTK", "ptr", objPtr)
	}
	callback.DeleteCallbacks(objPtr)
	Untrack(objPtr)
})

type Trackable interface {
	GoPointer() uintptr
	Unref()
	WeakRef(NotifyVar *gobject.WeakNotify, DataVar uintptr)
}

func SetFinalizer[T Trackable](tag string, widgetToBeFinalized T) {
	ptr := widgetToBeFinalized.GoPointer()
	Track(ptr, tag)
	if shouldLogLifecycle {
		logger.Debug("now managing cleanup for object", "ptr", ptr, "tag", tag)
	}
	widgetToBeFinalized.WeakRef(&onUnref, 0)
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
