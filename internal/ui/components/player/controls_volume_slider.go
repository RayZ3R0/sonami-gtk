package player

import (
	"math"

	"codeberg.org/dergs/tidalwave/internal/player"
	"codeberg.org/dergs/tidalwave/internal/signals"
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var volumeState = state.New(1.0)

func init() {
	player.OnVolumeChanged.On(func(volume float64) bool {
		// Cube root the volume to account for the logarithmic nature of human volume perception
		volumeState.SetValue(math.Pow(volume, 1.0/3.0))
		return signals.Continue
	})
}

func controlsVolumeSlider() schwifty.Scale {
	return Scale(gtk.OrientationVerticalValue).
		Inverted(true).
		Range(0, 1).
		Value(0.5).
		MinHeight(100).
		HExpand(true).
		CSS(`scale:active { background-color: transparent; }`).
		ConnectChangeValue(func(r gtk.Range, st gtk.ScrollType, value float64) bool {
			// Cube the value to account for the logarithmic nature of human volume perception
			player.SetVolume(math.Pow(value, 3))
			return false
		})
}
