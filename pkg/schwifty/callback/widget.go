package callback

import "github.com/jwijenbergh/puregotk/v4/gtk"

var (
	DestroyCallback = func(widget gtk.Widget) {
		CallbackHandler[any](widget.Object, "destroy", widget)
	}
	MapCallback = func(widget gtk.Widget) {
		CallbackHandler[any](widget.Object, "map", widget)
	}
	RealizedCallback = func(widget gtk.Widget) {
		CallbackHandler[any](widget.Object, "realize", widget)
	}
	UnmapCallback = func(widget gtk.Widget) {
		CallbackHandler[any](widget.Object, "unmap", widget)
	}
	UnrealizedCallback = func(widget gtk.Widget) {
		CallbackHandler[any](widget.Object, "unrealize", widget)
	}
)
