package tracking

import (
	"reflect"
	"runtime"

	"codeberg.org/puregotk/puregotk/v4/gobject"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/callback"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/utils/weak"
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
	Ref() *gobject.Object
	Unref()
	WeakRef(NotifyVar *gobject.WeakNotify, DataVar uintptr)
}

func SetFinalizer(tag string, widgetToBeFinalized Trackable) {
	widgetToBeFinalized.Ref()
	defer widgetToBeFinalized.Unref()

	Track(widgetToBeFinalized.GoPointer(), tag)
	widgetToBeFinalized.WeakRef(&onUnref, 0)

	ref := weak.NewObjectRef(widgetToBeFinalized)
	runtime.AddCleanup[any]((*any)(reflect.ValueOf(widgetToBeFinalized).UnsafePointer()), func(ref weak.ObjectRef) {
		callback.ScheduleOnMainThreadOncePure(func() {
			ref.Use(func(obj *gobject.Object) {
				if shouldLogLifecycle {
					logger.Debug("releasing reference on schwifty-managed object, GTK may still hold a reference to this object", "tag", tag, "ptr", obj.GoPointer())
				}
				obj.Unref()
				TrackGC(obj.Ptr)
			})
		})
	}, ref)

	if shouldLogLifecycle {
		logger.Debug("now managing cleanup for object", "ptr", widgetToBeFinalized.GoPointer(), "tag", tag)
	}
}

func SetLogLifecycle(enabled bool) {
	shouldLogLifecycle = enabled
}
