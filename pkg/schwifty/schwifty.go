package schwifty

import (
	"log/slog"

	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/callback"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/tracking"
)

var logger = slog.With("library", "schwifty")

var (
	shouldLogLifecycle = false
)

func OnMainThread(cb callback.MainThreadCallback, param uintptr) uint32 {
	return callback.OnMainThread(cb, param)
}

func OnMainThreadOnce(cb func(u uintptr), param uintptr) uint32 {
	return callback.OnMainThreadOnce(cb, param)
}

func OnMainThreadOncePure(cb func()) uint32 {
	return OnMainThreadOnce(func(uintptr) { cb() }, 0)
}

func SetLogLifecycle(enabled bool) {
	shouldLogLifecycle = enabled
	callback.SetLogLifecycle(enabled)
	tracking.SetLogLifecycle(enabled)
}
