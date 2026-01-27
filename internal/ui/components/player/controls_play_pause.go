package player

import (
	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/player"
	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var (
	playPauseIconState = state.NewStateful("play-symbolic")
)

func init() {
	player.PlaybackStateChanged.On(func(state *player.PlaybackState) bool {
		schwifty.OnMainThreadOncePure(func() {
			switch state.Status {
			case player.PlaybackStatusPlaying:
				playPauseIconState.SetValue("pause-symbolic")
			case player.PlaybackStatusPaused, player.PlaybackStatusStopped:
				playPauseIconState.SetValue("play-symbolic")
			}
		})

		return signals.Continue
	})
}

func controlsPlayPause() schwifty.Button {
	var controllableStateSub *signals.Subscription
	return Button().
		TooltipText(gettext.Get("Play / Pause")).
		IconName(playPauseIconState.Value()).
		BindIconName(playPauseIconState).
		BindSensitive(isControllable).
		CornerRadius(21).
		HPadding(32).
		VPadding(9).
		CSS(`button:not(:hover) { background-color: var(--accent-bg-color); color: var(--accent-fg-color); }`).
		ConnectClicked(func(b gtk.Button) {
			player.PlayPause()
		}).
		ConnectConstruct(func(b *gtk.Button) {
			ptr := b.GoPointer()
			player.PlaybackStateChanged.On(func(ps *player.PlaybackState) bool {
				b := gtk.ButtonNewFromInternalPtr(ptr)
				schwifty.OnMainThreadOncePure(func() {
					if ps.Loading {
						child := Spinner().ToGTK()
						b.SetChild(child)
					}
				})
				return signals.Continue
			})
		}).
		ConnectDestroy(func(w gtk.Widget) {
			player.PlaybackStateChanged.Unsubscribe(controllableStateSub)
		})
}
