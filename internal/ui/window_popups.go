package ui

import (
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/internal/secrets"
	"github.com/RayZ3R0/sonami-gtk/internal/settings"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	"codeberg.org/puregotk/puregotk/v4/adw"
	"codeberg.org/puregotk/puregotk/v4/gtk"
)

func (w *Window) PresentSecretServiceError(err *secrets.ServiceError) {
	if settings.General().ShouldHideSecretServiceWarning() {
		return
	}

	// ConnectResponse is broken with puregotk, so we have to manually hack our way
	AlertDialog(err.Title, err.Body).
		WithCSSClass("no-response").
		ConnectConstruct(func(ad *adw.AlertDialog) {
			checkbox := gtk.NewCheckButtonWithLabel(gettext.Get("Don't show again"))
			checkbox.SetHalign(gtk.AlignBaselineCenterValue)
			checkbox.AddCssClass("space-2")

			ad.SetExtraChild(
				VStack(
					checkbox,
					Button().Label(gettext.Get("Continue")).WithCSSClass("destructive-action").VPadding(10).ConnectClicked(func(b gtk.Button) {
						if checkbox.GetActive() {
							settings.General().SetHideSecretServiceWarning(true)
						}
						ad.Close()
					}),
				).Spacing(12).ToGTK(),
			)
			ad.Present(&w.Widget)
		})()
}
