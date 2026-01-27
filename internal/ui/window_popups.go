package ui

import (
	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/secrets"
	"codeberg.org/dergs/tonearm/internal/settings"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func (w *Window) PresentSecretServiceError(err *secrets.ServiceError) {
	if settings.General().ShouldHideSecretServiceWarning() {
		return
	}

	// ConnectResponse is broken with puregotk, so we have to manually hack our way
	AlertDialog(err.Title, err.Body).
		WithCSSClass("no-response").
		ConnectConstruct(func(ad *adw.AlertDialog) {
			ad.SetExtraChild(
				VStack(
					Button().Label(gettext.Get("Continue")).VPadding(10).ConnectClicked(func(b gtk.Button) {
						ad.Close()
					}),
					Button().Label(gettext.Get("Never Show This Again")).WithCSSClass("destructive-action").VPadding(10).ConnectClicked(func(b gtk.Button) {
						settings.General().SetHideSecretServiceWarning(true)
						ad.Close()
					}),
				).Spacing(12).ToGTK(),
			)
			ad.Present(&w.Widget)
		})()
}
