package settings

import (
	"fmt"
	"log/slog"
	"sync"

	"codeberg.org/dergs/tonearm/internal/g"
	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/pkg/schwifty/callback"
	v1 "codeberg.org/dergs/tonearm/pkg/tidalapi/models/v1"
	"codeberg.org/dergs/tonearm/pkg/utils/cutil"
	"github.com/jwijenbergh/puregotk/v4/gio"
	"github.com/jwijenbergh/puregotk/v4/glib"
	"github.com/jwijenbergh/puregotk/v4/gobject"
)

type PropertyChangedCallback func(*glib.Variant) bool

type PlayerSettings struct {
	settings *gio.Settings

	changedCallback     uint32
	changedCallbackLock sync.RWMutex
}

func (p *PlayerSettings) initCallbackAggregator() {
	p.changedCallbackLock.Lock()
	defer p.changedCallbackLock.Unlock()

	p.changedCallback = p.settings.ConnectChanged(g.Ptr(func(s gio.Settings, setting string) {
		setting = cutil.ParseNullTerminatedString(setting)
		if setting == "" {
			return
		}

		val := s.GetValue(setting)

		callback.CallbackHandler[any](s.Object, fmt.Sprintf("changed::%s", setting), val)
	}))
}

func (p *PlayerSettings) ConnectAudioQualityChanged(cb func(v1.AudioQuality) bool) {
	if p.changedCallback == 0 {
		p.initCallbackAggregator()
	}

	if !callback.HasCallback(p.settings.Object, "changed::audio-quality") {
		// Ensure that audio-quality gets watched by gio,
		// since ConnectChanged specifies "
		// 	Note that @settings only emits this signal if you have read
		// 	@key at least once while a signal handler was already connected for @key.
		// "
		p.settings.GetString("audio-quality")
	}

	subscription := new(int)
	(*subscription) = callback.HandleCallback(p.settings.Object, "changed::audio-quality", func(val *glib.Variant) {
		var quality v1.AudioQuality
		switch val.GetString(nil) {
		case string(v1.AudioQualityHighResLossless):
			quality = v1.AudioQualityHighResLossless
		case string(v1.AudioQualityLossless):
			quality = v1.AudioQualityLossless
		case string(v1.AudioQualityHighRes):
			quality = v1.AudioQualityHighRes
		case string(v1.AudioQualityLossy):
			quality = v1.AudioQualityLossy
		default:
			slog.Error("configured audio quality was invalid, defaulting to maximum", "configured", quality)
			quality = v1.AudioQualityHighResLossless
		}

		next := cb(quality)
		if next == signals.Unsubscribe {
			callback.DeleteCallback(p.settings.Object, "changed::audio-quality", *subscription)
		}

		return
	})
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
