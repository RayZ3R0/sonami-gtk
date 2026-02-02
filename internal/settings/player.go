package settings

import (
	"log/slog"
	"maps"
	"sync"

	"codeberg.org/dergs/tonearm/internal/g"
	"codeberg.org/dergs/tonearm/internal/signals"
	v1 "codeberg.org/dergs/tonearm/pkg/tidalapi/models/v1"
	"codeberg.org/dergs/tonearm/pkg/utils/cutil"
	"github.com/google/uuid"
	"github.com/jwijenbergh/puregotk/v4/gio"
	"github.com/jwijenbergh/puregotk/v4/glib"
	"github.com/jwijenbergh/puregotk/v4/gobject"
)

type PropertyChangedCallback func(*glib.Variant) bool

type PlayerSettings struct {
	settings *gio.Settings

	cbMutex   sync.RWMutex
	callbacks map[string]map[*signals.Subscription]PropertyChangedCallback
}

func (p *PlayerSettings) initCallbackAggregator() {
	p.cbMutex.Lock()
	defer p.cbMutex.Unlock()

	p.callbacks = make(map[string]map[*signals.Subscription]PropertyChangedCallback)

	p.settings.ConnectChanged(g.Ptr(func(s gio.Settings, setting string) {
		setting = cutil.ParseNullTerminatedString(setting)
		if setting == "" {
			return
		}

		val := s.GetValue(setting)

		if callbacks, ok := p.callbacks[setting]; ok {
			p.cbMutex.RLock()
			callbacks := maps.Clone(callbacks)
			p.cbMutex.RUnlock()

			for id, callback := range callbacks {
				if callback(val) == signals.Unsubscribe {
					p.cbMutex.Lock()
					delete(p.callbacks[setting], id)
					p.cbMutex.Unlock()
				}
			}
		}
	}))
}

func (p *PlayerSettings) ConnectAudioQualityChanged(cb func(v1.AudioQuality) bool) {
	var array map[*signals.Subscription]PropertyChangedCallback
	var ok bool

	if p.callbacks == nil {
		p.initCallbackAggregator()
	}

	p.settings.GetString("audio-quality")

	if array, ok = p.callbacks["audio-quality"]; !ok {
		array = make(map[*signals.Subscription]PropertyChangedCallback, 0)
	}
	id := signals.Subscription(uuid.New())
	array[&id] = func(val *glib.Variant) bool {
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

		return cb(quality)
	}

	p.callbacks["audio-quality"] = array
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
