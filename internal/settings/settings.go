package settings

import (
	"runtime"

	"codeberg.org/dergs/tonearm/internal/g"
	"codeberg.org/dergs/tonearm/pkg/schwifty/callback"
	"codeberg.org/dergs/tonearm/pkg/schwifty/tracking"
	"github.com/jwijenbergh/puregotk/v4/gio"
)

//go:generate glib-compile-schemas .

var General = g.Lazy(func() *GeneralSettings {
	return &GeneralSettings{
		finalize(gio.NewSettings("dev.dergs.Tonearm")),
	}
})

var Playback = g.Lazy(func() *PlaybackSettings {
	return &PlaybackSettings{
		finalize(gio.NewSettings("dev.dergs.Tonearm.playback")),
	}
})

var Performance = g.Lazy(func() *PerformanceSettings {
	return &PerformanceSettings{
		finalize(gio.NewSettings("dev.dergs.Tonearm.performance")),
	}
})

var Player = g.Lazy(func() *PlayerSettings {
	settings := gio.NewSettings("dev.dergs.Tonearm.player")
	settings.ConnectChanged(&callback.GioSettingsChangedCallback)
	tracking.Track(settings.GoPointer(), "Settings")
	return &PlayerSettings{
		settings: settings,
	}
})

var Scrobbling = g.Lazy(func() *ScrobblingSettings {
	return &ScrobblingSettings{
		finalize(gio.NewSettings("dev.dergs.Tonearm.scrobbling")),
	}
})

func finalize(settings *gio.Settings) *gio.Settings {
	runtime.SetFinalizer(settings, func(s *gio.Settings) {
		s.Unref()
	})
	return settings
}
