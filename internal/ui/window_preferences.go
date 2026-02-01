package ui

import (
	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/settings"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var preferencesGeneral = PreferencesPage(
	PreferencesGroup(
		SwitchRow().
			Title(gettext.Get("Allow Background Activity")).
			Subtitle(gettext.Get("Allow Tonearm to run in the background by hiding the player window instead of quitting the application")).
			ConnectConstruct(func(sr *adw.SwitchRow) {
				settings.General().BindRunInBackground(&sr.Object, "active")
			}),
	).
		Title(gettext.Get("Background Activity")).
		Description(gettext.Get("Configure the behaviour of Tonearm when running in the background.")),
	PreferencesGroup(
		EntryRow().
			Title(gettext.Get("Default Page")).
			ConnectConstruct(func(sr *adw.EntryRow) {
				settings.General().BindDefaultPage(&sr.Object, "text")
			}),
		SpinRow(
			gtk.NewAdjustment(0, 1, 100, 1, 1, 1),
			1,
			0,
		).Title(gettext.Get("History Length")).
			Subtitle(gettext.Get("Maximum history length before dropping old entries.")).
			ConnectConstruct(func(sr *adw.SpinRow) {
				settings.Performance().BindMaxRouterHistorySize(&sr.Object, "value")
			}),
	).
		Title(gettext.Get("Navigation Behaviour")).
		Description(gettext.Get("Configure the behaviour of Tonearm when navigating between pages.")),
).Title(gettext.Get("General")).IconName("settings-symbolic")

var preferencesPlayback = PreferencesPage(
	PreferencesGroup(
		SwitchRow().
			Title(gettext.Get("Enable autoplay")).
			Subtitle(gettext.Get("Allow Tonearm to start a mix with the current playing song when you're at the end of an album/queue")).
			ConnectConstruct(func(sr *adw.SwitchRow) {
				settings.Playback().BindAllowAutoplay(&sr.Object, "active")
			}),
	).Title(gettext.Get("Playback")).Description(gettext.Get("Configure the behaviour of Tonearm regarding playback.")),
).Title(gettext.Get("Playback")).IconName("media-playback-start-symbolic")

var preferencesPerformance = PreferencesPage(
	PreferencesGroup(
		SwitchRow().
			Title(gettext.Get("Allow Media Card Images")).
			Subtitle(gettext.Get("Allow Tonearm to load images for media card buttons.")).
			ConnectConstruct(func(sr *adw.SwitchRow) {
				settings.Performance().BindAllowMediaCardImages(&sr.Object, "active")
			}),
		SwitchRow().
			Title(gettext.Get("Allow Shortcut Images")).
			Subtitle(gettext.Get("Allow Tonearm to load images for shortcut buttons.")).
			ConnectConstruct(func(sr *adw.SwitchRow) {
				settings.Performance().BindAllowShortcutImages(&sr.Object, "active")
			}),
		SwitchRow().
			Title(gettext.Get("Allow Tracklist Images")).
			Subtitle(gettext.Get("Allow Tonearm to load images for tracklists with the cover column.")).
			ConnectConstruct(func(sr *adw.SwitchRow) {
				settings.Performance().BindAllowTracklistImages(&sr.Object, "active")
			}),
		SwitchRow().
			Title(gettext.Get("Cache Images")).
			Subtitle(gettext.Get("Allow Tonearm to temporarily store images on the file system to improve performance and reduce network traffic.")).
			ConnectConstruct(func(sr *adw.SwitchRow) {
				settings.Performance().BindCacheImages(&sr.Object, "active")
			}),
	).Title(gettext.Get("Images")).Description(gettext.Get("Configure the behaviour of Tonearm regarding images.")),
).Title(gettext.Get("Performance")).IconName("speedometer5-symbolic")

var preferencesScrobbling = PreferencesPage(
	PreferencesGroup(
		SwitchRow().
			Title(gettext.Get("Enable ListenBrainz")).
			Subtitle(gettext.Get("Allow Tonearm to send scrobbling data to ListenBrainz")).
			ConnectConstruct(func(sr *adw.SwitchRow) {
				settings.Scrobbling().BindEnableListenBrainz(&sr.Object, "active")
			}),
		PasswordEntryRow().
			Title(gettext.Get("ListenBrainz API Token")).
			ConnectConstruct(func(sr *adw.PasswordEntryRow) {
				settings.Scrobbling().BindListenBrainzToken(&sr.Object, "text")
			}),
		EntryRow().
			Title(gettext.Get("ListenBrainz API URL")).
			ConnectConstruct(func(sr *adw.EntryRow) {
				settings.Scrobbling().BindListenBrainzUrl(&sr.Object, "text")
			}),
	).
		Title(gettext.Get("ListenBrainz")).
		Description(gettext.Get("Configure Tonearm to send scrobbling data to ListenBrainz.")),
).Title(gettext.Get("Scrobbling")).IconName("podcast-symbolic")

func (w *Window) PresentPreferences() {
	PreferencesDialog(
		preferencesGeneral,
		preferencesPlayback,
		preferencesPerformance,
		preferencesScrobbling,
	).Present(w)
}
