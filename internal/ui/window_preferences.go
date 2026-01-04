package ui

import (
	"codeberg.org/dergs/tidalwave/internal/settings"
	"github.com/jwijenbergh/puregotk/v4/adw"
)

func (w *Window) PresentPreferences() {
	preferences := adw.NewPreferencesDialog()

	allowBg := adw.NewSwitchRow()
	allowBg.SetTitle("Allow Background Activity")
	allowBg.SetSubtitle("Allow Tidal Wave to run in the background")
	settings.General().BindRunInBackground(&allowBg.Object, "active")

	bgPreferences := adw.NewPreferencesGroup()
	bgPreferences.SetTitle("Background Activity")
	bgPreferences.SetDescription("Configure the behaviour of Tidal Wave when running in the background.")
	bgPreferences.Add(&allowBg.Widget)
	allowBg.Unref()

	mainPage := adw.NewPreferencesPage()
	mainPage.SetTitle("General")
	mainPage.Add(bgPreferences)
	bgPreferences.Unref()

	preferences.Add(mainPage)
	mainPage.Unref()

	preferences.Present(&w.Widget)
	preferences.Unref()
}
