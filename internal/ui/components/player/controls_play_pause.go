package player

import (
	"codeberg.org/dergs/tidalwave/internal/player"
	"codeberg.org/dergs/tidalwave/internal/signals"
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var (
	playPauseIconState      = state.NewStateful("media-playback-start-symbolic")
	playPauseSensitiveState = state.NewStateful(true)
)

func init() {
	player.OnStateChanged.On(func(state player.State) bool {
		switch state.Status {
		case player.StatusPlaying:
			playPauseIconState.SetValue("media-playback-pause-symbolic")
			playPauseSensitiveState.SetValue(true)
		case player.StatusPaused, player.StatusStopped:
			playPauseIconState.SetValue("media-playback-start-symbolic")
			playPauseSensitiveState.SetValue(true)
		default:
			playPauseIconState.SetValue("media-playback-start-symbolic")
			playPauseSensitiveState.SetValue(false)
		}

		return signals.Continue
	})
}

func controlsPlayPause() schwifty.Button {
	return Button().
		IconName(playPauseIconState.Value()).
		BindIconName(playPauseIconState).
		BindSensitive(playPauseSensitiveState).
		CornerRadius(21).
		HPadding(32).
		VPadding(9).
		CSS(`button:not(:hover) { background-color: var(--accent-bg-color); }`).
		ConnectClicked(func(b gtk.Button) {
			player.PlayPause()
		})
}
