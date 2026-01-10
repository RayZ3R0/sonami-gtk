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
		).Title("General").IconName("settings-symbolic"),
		PreferencesPage(
			PreferencesGroup(
				SwitchRow().
					Title("Allow Media Card Images").
					Subtitle("Allow Tidal Wave to load images for media card buttons.").
					ConnectConstruct(func(sr *adw.SwitchRow) {
						settings.Performance().BindAllowMediaCardImages(&sr.Object, "active")
					}),
				SwitchRow().
					Title("Allow Shortcut Images").
					Subtitle("Allow Tidal Wave to load images for shortcut buttons.").
					ConnectConstruct(func(sr *adw.SwitchRow) {
						settings.Performance().BindAllowShortcutImages(&sr.Object, "active")
					}),
				SwitchRow().
					Title("Allow Tracklist Images").
					Subtitle("Allow Tidal Wave to load images for tracklists with the cover column.").
					ConnectConstruct(func(sr *adw.SwitchRow) {
						settings.Performance().BindAllowTracklistImages(&sr.Object, "active")
					}),
				SwitchRow().
					Title("Cache Images").
					Subtitle("Allow Tidal Wave to temporarily store images on the file system to improve performance and reduce network traffic.").
					ConnectConstruct(func(sr *adw.SwitchRow) {
						settings.Performance().BindCacheImages(&sr.Object, "active")
					}),
			).
				Title("Images").
				Description("Configure the behaviour of Tidal Wave regarding images."),
		).Title("Performance").IconName("power-profile-performance-symbolic"),
	).Present(w)
}
