package ui

import (
	"codeberg.org/dergs/tonearm/internal/settings"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var preferencesGeneral = PreferencesPage(
	PreferencesGroup(
		SwitchRow().
			Title("Allow Background Activity").
			Subtitle("Allow Tonearm to run in the background by hiding the player window instead of quitting the application").
			ConnectConstruct(func(sr *adw.SwitchRow) {
				settings.General().BindRunInBackground(&sr.Object, "active")
			}),
	).
		Title("Background Activity").
		Description("Configure the behaviour of Tonearm when running in the background."),
	PreferencesGroup(
		EntryRow().
			Title("Default Page").
			ConnectConstruct(func(sr *adw.EntryRow) {
				settings.General().BindDefaultPage(&sr.Object, "text")
			}),
		SpinRow(
			gtk.NewAdjustment(0, 1, 100, 1, 1, 1),
			1,
			0,
		).Title("History Length").Subtitle("Maximum history length before dropping old entries.").ConnectConstruct(func(sr *adw.SpinRow) {
			settings.Performance().BindMaxRouterHistorySize(&sr.Object, "value")
		}),
	).
		Title("Navigation Behaviour").
		Description("Configure the behaviour of Tonearm when navigating between pages."),
).Title("General").IconName("settings-symbolic")

var preferencesPerformance = PreferencesPage(
	PreferencesGroup(
		SwitchRow().
			Title("Allow Media Card Images").
			Subtitle("Allow Tonearm to load images for media card buttons.").
			ConnectConstruct(func(sr *adw.SwitchRow) {
				settings.Performance().BindAllowMediaCardImages(&sr.Object, "active")
			}),
		SwitchRow().
			Title("Allow Shortcut Images").
			Subtitle("Allow Tonearm to load images for shortcut buttons.").
			ConnectConstruct(func(sr *adw.SwitchRow) {
				settings.Performance().BindAllowShortcutImages(&sr.Object, "active")
			}),
		SwitchRow().
			Title("Allow Tracklist Images").
			Subtitle("Allow Tonearm to load images for tracklists with the cover column.").
			ConnectConstruct(func(sr *adw.SwitchRow) {
				settings.Performance().BindAllowTracklistImages(&sr.Object, "active")
			}),
		SwitchRow().
			Title("Cache Images").
			Subtitle("Allow Tonearm to temporarily store images on the file system to improve performance and reduce network traffic.").
			ConnectConstruct(func(sr *adw.SwitchRow) {
				settings.Performance().BindCacheImages(&sr.Object, "active")
			}),
	).Title("Images").Description("Configure the behaviour of Tonearm regarding images."),
).Title("Performance").IconName("speedometer5-symbolic")

var preferencesScrobbling = PreferencesPage(
	PreferencesGroup(
		SwitchRow().
			Title("Enable ListenBrainz").
			Subtitle("Allow Tonearm to send scrobbling data to ListenBrainz").
			ConnectConstruct(func(sr *adw.SwitchRow) {
				settings.Scrobbling().BindEnableListenBrainz(&sr.Object, "active")
			}),
		PasswordEntryRow().
			Title("ListenBrainz API Token").
			ConnectConstruct(func(sr *adw.PasswordEntryRow) {
				settings.Scrobbling().BindListenBrainzToken(&sr.Object, "text")
			}),
	).
		Title("ListenBrainz").
		Description("Configure Tonearm to send scrobbling data to ListenBrainz."),
).Title("Scrobbling").IconName("podcast-symbolic")

func (w *Window) PresentPreferences() {
	PreferencesDialog(
		preferencesGeneral,
		preferencesPerformance,
		preferencesScrobbling,
	).Present(w)
}
