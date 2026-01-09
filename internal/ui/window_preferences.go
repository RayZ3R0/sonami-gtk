package ui

import (
	"codeberg.org/dergs/tidalwave/internal/settings"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"github.com/jwijenbergh/puregotk/v4/adw"
)

func (w *Window) PresentPreferences() {
	PreferencesDialog(
		PreferencesPage(
			PreferencesGroup(
				SwitchRow().
					Title("Allow Background Activity").
					Subtitle("Allow Tidal Wave to run in the background by hiding the player window instead of quitting the application").
					ConnectConstruct(func(sr *adw.SwitchRow) {
						settings.General().BindRunInBackground(&sr.Object, "active")
					}),
			).
				Title("Background Activity").
				Description("Configure the behaviour of Tidal Wave when running in the background."),
		).Title("General"),
	).Present(w)
}
