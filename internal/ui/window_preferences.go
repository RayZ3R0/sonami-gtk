package ui

import (
	"codeberg.org/dergs/tidalwave/internal/settings"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"github.com/jwijenbergh/puregotk/v4/adw"
)

func (w *Window) PresentPreferences() {
	allowBg := adw.NewSwitchRow()
	allowBg.SetTitle("Allow Background Activity")
	allowBg.SetSubtitle("Allow Tidal Wave to run in the background")
	settings.General().BindRunInBackground(&allowBg.Object, "active")

	PreferencesDialog(
		PreferencesPage(
			PreferencesGroup(
				ManagedWidget(&allowBg.Widget),
			).
				Title("Background Activity").
				Description("Configure the behaviour of Tidal Wave when running in the background."),
		).Title("General"),
	).Present(w)
}
