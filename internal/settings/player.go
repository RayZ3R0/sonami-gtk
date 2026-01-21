package settings

import (
	"log/slog"

	v1 "codeberg.org/dergs/tonearm/pkg/tidalapi/models/v1"
	"github.com/jwijenbergh/puregotk/v4/gio"
	"github.com/jwijenbergh/puregotk/v4/gobject"
)

type PlayerSettings struct {
	settings *gio.Settings
}

func (p *PlayerSettings) GetAudioQuality() v1.AudioQuality {
	quality := p.settings.GetString("audio-quality")
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

func (p *PlayerSettings) BindVolume(target *gobject.Object, property string) {
	p.settings.Bind("volume", target, property, gio.GSettingsBindNoSensitivityValue)
}

func (p *PlayerSettings) GetVolume() float64 {
	return p.settings.GetDouble("volume")
}

func (p *PlayerSettings) SetVolume(volume float64) {
	p.settings.SetDouble("volume", volume)
}

func (p *PlayerSettings) SetAudioQuality(quality v1.AudioQuality) {
	p.settings.SetString("audio-quality", string(quality))
}
