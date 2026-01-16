package callback

import "github.com/jwijenbergh/puregotk/v4/adw"

var (
	AlertDialogCloseAttempt = func(widget adw.Dialog) {
		CallbackHandler[any](widget.Object, "close-attempt", widget)
	}
)
