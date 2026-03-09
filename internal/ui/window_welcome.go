package ui

import (
	"codeberg.org/puregotk/puregotk/v4/adw"
	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/internal/settings"
)

// MaybePresentWelcome shows a welcome dialog if the streaming instances URL
// has not been configured yet. It gently asks the user to paste the URL and
// tells them where to find the setting later (Preferences → Streaming).
// Returns true if the dialog was shown.
func (w *Window) MaybePresentWelcome() bool {
	if settings.Streaming().GetInstancesURL() != "" {
		return false
	}

	dialog := adw.NewAlertDialog(
		gettext.Get("Welcome to Sonami"),
		gettext.Get("To enable music playback, you need to provide a streaming instances URL.\n\nPaste the URL of your instances JSON below to get started. You can change this later in Preferences → Streaming."),
	)

	// Add an entry for the URL as extra content
	entry := gtk.NewEntry()
	entry.SetPlaceholderText("https://...")
	entry.SetHexpand(true)

	clamp := adw.NewClamp()
	clamp.SetMaximumSize(500)
	clamp.SetChild(&entry.Widget)

	dialog.SetExtraChild(&clamp.Widget)

	dialog.AddResponse("skip", gettext.Get("Skip for Now"))
	dialog.AddResponse("save", gettext.Get("Save & Continue"))
	dialog.SetResponseAppearance("save", adw.ResponseSuggestedValue)
	dialog.SetDefaultResponse("save")
	dialog.SetCloseResponse("skip")

	dialog.ConnectResponse(new(func(_ adw.AlertDialog, response string) {
		if response == "save" {
			text := entry.GetText()
			if text != "" {
				settings.Streaming().SetInstancesURL(text)
			}
		}
	}))

	dialog.Present(&w.Widget)
	return true
}
