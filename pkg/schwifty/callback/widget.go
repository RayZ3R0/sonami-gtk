package callback

import "github.com/jwijenbergh/puregotk/v4/gtk"

var (
	RealizedCallback = func(widget gtk.Widget) {
		CallbackHandler[any](widget, "realize", widget)
	}
	UnrealizedCallback = func(widget gtk.Widget) {
		CallbackHandler[any](widget, "unrealize", widget)
	}
)
