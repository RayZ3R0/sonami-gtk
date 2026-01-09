package callback

import "github.com/jwijenbergh/puregotk/v4/gtk"

var (
	ButtonClickedCallback = func(widget gtk.Button) {
		CallbackHandler[any](widget.Object, "clicked", widget)
	}
)
