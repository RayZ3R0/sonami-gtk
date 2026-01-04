package settings

import (
	"log/slog"

	v1 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v1"
	"github.com/jwijenbergh/puregotk/v4/gio"
	"github.com/jwijenbergh/puregotk/v4/gobject"
)

type Player struct {
	settings *gio.Settings
}

func (p *Player) GetAudioQuality() v1.AudioQuality {
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

func (p *Player) BindVolume(target *gobject.Object, property string) {
	p.settings.Bind("volume", target, property, gio.GSettingsBindNoSensitivityValue)
}

func (p *Player) SetAudioQuality(quality v1.AudioQuality) {
	p.settings.SetString("audio-quality", string(quality))
}
