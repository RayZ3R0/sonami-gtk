package settings

import (
	"log/slog"

	"github.com/RayZ3R0/sonami-gtk/internal/signals"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/callback"
	v1 "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/models/v1"
	"codeberg.org/puregotk/puregotk/v4/gio"
	"codeberg.org/puregotk/puregotk/v4/gobject"
)

type PlayerSettings struct {
	settings *gio.Settings
}

func (p *PlayerSettings) ConnectAudioQualityChanged(cb func(v1.AudioQuality) bool) {
	// Ensure that audio-quality gets watched by gio,
	// since ConnectChanged specifies "
	// 	Note that @settings only emits this signal if you have read
	// 	@key at least once while a signal handler was already connected for @key.
	// "
	if cb(p.GetAudioQuality()) == signals.Unsubscribe {
		return
	}

	var callbackId int
	callbackId = callback.HandleCallback(p.settings.Object, "changed", func(settings gio.Settings, setting string) {
		if setting == "audio-quality" {
			if cb(p.GetAudioQuality()) == signals.Unsubscribe {
				callback.DeleteCallback(p.settings.Object, "changed", callbackId)
			}
		}
	})
}

func (p *PlayerSettings) GetAudioQuality() v1.AudioQuality {
	return p.parseAudioQuality(p.settings.GetString("audio-quality"))
}

func (p *PlayerSettings) SetAudioQuality(quality v1.AudioQuality) {
	p.settings.SetString("audio-quality", string(quality))
}

func (p *PlayerSettings) BindVolume(target *gobject.Object, property string) {
	p.settings.Bind("volume", target, property, gio.GSettingsBindNoSensitivityValue)
}

func (p *PlayerSettings) GetVolume() float64 {
	return p.settings.GetDouble("volume")
}

func (p *PlayerSettings) SetVolume(volume float64) {
	p.settings.SetDouble("volume", volume)
}

func (p *PlayerSettings) parseAudioQuality(quality string) v1.AudioQuality {
	switch quality {
	case string(v1.AudioQualityHighResLossless):
		return v1.AudioQualityHighResLossless
	case string(v1.AudioQualityLossless):
		return v1.AudioQualityLossless
	case string(v1.AudioQualityHighRes):
		return v1.AudioQualityHighRes
	case string(v1.AudioQualityLossy):
		return v1.AudioQualityLossy
	default:
		slog.Error("configured audio quality was invalid, defaulting to maximum", "configured", quality)
		return v1.AudioQualityHighResLossless
	}
}
