package callback

import (
	"codeberg.org/puregotk/puregotk/v4/adw"
	"github.com/RayZ3R0/sonami-gtk/pkg/utils/cutil"
)

var (
	AlertDialogClosed = func(widget adw.Dialog) {
		CallbackHandler[any](widget.Object, "closed", widget)
	}
	AlertDialogCloseAttempt = func(widget adw.Dialog) {
		CallbackHandler[any](widget.Object, "close-attempt", widget)
	}
	AlertDialogResponse = func(widget adw.AlertDialog, response string) {
		response = cutil.ParseNullTerminatedString(response)
		CallbackHandler[string](widget.Object, "response", widget, response)
	}
)

var (
	ActionRowActivated = func(widget adw.ActionRow) {
		CallbackHandler[any](widget.Object, "activated", widget)
	}
)
