package settings

import (
	"runtime"

	"codeberg.org/dergs/tonearm/internal/g"
	"github.com/jwijenbergh/puregotk/v4/gio"
)

//go:generate glib-compile-schemas .

var General = g.Lazy(func() *GeneralSettings {
	return &GeneralSettings{
		finalize(gio.NewSettings("dev.dergs.Tonearm")),
	}
})

var Performance = g.Lazy(func() *PerformanceSettings {
	return &PerformanceSettings{
		finalize(gio.NewSettings("dev.dergs.Tonearm.performance")),
	}
})

var Player = g.Lazy(func() *PlayerSettings {
	return &PlayerSettings{
		finalize(gio.NewSettings("dev.dergs.Tonearm.player")),
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
