package settings

import (
	"codeberg.org/puregotk/puregotk/v4/gio"
	"github.com/RayZ3R0/sonami-gtk/internal/g"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/callback"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/tracking"
	"github.com/RayZ3R0/sonami-gtk/pkg/utils/cutil"
)

//go:generate glib-compile-schemas .

var (
	GioSettingsChangedCallback = func(settings gio.Settings, setting string) {
		callback.CallbackHandler[any](settings.Object, "changed", settings, cutil.ParseNullTerminatedString(setting))
	}
)

var General = g.Lazy(func() *GeneralSettings {
	return &GeneralSettings{
		finalize(gio.NewSettings("io.github.rayz3r0.SonamiGtk")),
	}
})

var Playback = g.Lazy(func() *PlaybackSettings {
	settings := gio.NewSettings("io.github.rayz3r0.SonamiGtk.playback")
	settings.ConnectChanged(&GioSettingsChangedCallback)
	tracking.Track(settings.GoPointer(), "Settings")
	return &PlaybackSettings{
		finalize(settings),
	}
})

var Performance = g.Lazy(func() *PerformanceSettings {
	return &PerformanceSettings{
		finalize(gio.NewSettings("io.github.rayz3r0.SonamiGtk.performance")),
	}
})

var Player = g.Lazy(func() *PlayerSettings {
	settings := gio.NewSettings("io.github.rayz3r0.SonamiGtk.player")
	settings.ConnectChanged(&GioSettingsChangedCallback)
	tracking.Track(settings.GoPointer(), "Settings")
	return &PlayerSettings{
		finalize(settings),
	}
})

var Scrobbling = g.Lazy(func() *ScrobblingSettings {
	return &ScrobblingSettings{
		finalize(gio.NewSettings("io.github.rayz3r0.SonamiGtk.scrobbling")),
	}
})

var Discord = g.Lazy(func() *DiscordSettings {
	return &DiscordSettings{
		finalize(gio.NewSettings("io.github.rayz3r0.SonamiGtk.discord")),
	}
})

var Lyrics = g.Lazy(func() *LyricsSettings {
	return &LyricsSettings{
		finalize(gio.NewSettings("io.github.rayz3r0.SonamiGtk.lyrics")),
	}
})

var Streaming = g.Lazy(func() *StreamingSettings {
	return &StreamingSettings{
		finalize(gio.NewSettings("io.github.rayz3r0.SonamiGtk.streaming")),
	}
})

func finalize(settings *gio.Settings) *gio.Settings {
	tracking.SetFinalizer("Settings", settings)
	return settings
}
