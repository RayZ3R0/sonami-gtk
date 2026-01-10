package settings

import (
	"runtime"

	"codeberg.org/dergs/tidalwave/internal/g"
	"github.com/jwijenbergh/puregotk/v4/gio"
)

//go:generate glib-compile-schemas .

var General = g.Lazy(func() *GeneralSettings {
	return &GeneralSettings{
		finalize(gio.NewSettings("org.codeberg.dergs.tidalwave")),
	}
})

var Performance = g.Lazy(func() *PerformanceSettings {
	return &PerformanceSettings{
		finalize(gio.NewSettings("org.codeberg.dergs.tidalwave.performance")),
	}
})

func PlayerSettings() *Player {
	return &Player{
		finalize(gio.NewSettings("org.codeberg.dergs.tidalwave.player")),
	}
}

func finalize(settings *gio.Settings) *gio.Settings {
	runtime.SetFinalizer(settings, func(s *gio.Settings) {
		s.Unref()
	})
	return settings
}
