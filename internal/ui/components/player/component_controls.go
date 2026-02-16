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

const (
	repeatListIconName  = "playlist-repeat-symbolic"
	repeatTrackIconName = "playlist-repeat-song-symbolic"
	pauseIconName       = "pause-symbolic"
	playIconName        = "play-symbolic"
)

var controlButton = Button().
	BindSensitive(isControllableState).
	MinHeight(34).MinWidth(34).
	WithCSSClass("flat").
	VAlign(gtk.AlignCenterValue)

func controls() schwifty.Box {
	var (
		pauseIcon   = Image().FromIconName(pauseIconName)()
		playIcon    = Image().FromIconName(playIconName)()
		playSpinner = Spinner().SizeRequest(16, 16)()
	)

	var (
		playPauseChildState       = state.NewStateful[any](nil)
		actualPlayPauseChildState = state.NewStateful[any](nil)

		repeatClassState = state.NewStateful("")
		repeatIconState  = state.NewStateful(repeatListIconName)

		shuffleClassState = state.NewStateful("")
	)

	player.PlaybackStateChanged.On(func(state *player.PlaybackState) bool {
		schwifty.OnMainThreadOncePure(func() {
			var val any
			switch state.Status {
			case player.PlaybackStatusPlaying:
				val = pauseIcon
			case player.PlaybackStatusPaused, player.PlaybackStatusStopped:
				val = playIcon
			}

			if isControllableState.Value() {
				actualPlayPauseChildState.SetValue(val)
			} else {
				actualPlayPauseChildState.SetValue(playSpinner)
			}

			playPauseChildState.SetValue(val)
		})

		return signals.Continue
	})

	isControllableState.AddCallback(func(newValue bool) {
		schwifty.OnMainThreadOncePure(func() {
			if newValue {
				actualPlayPauseChildState.SetValue(playPauseChildState.Value())
			} else {
				actualPlayPauseChildState.SetValue(playSpinner)
			}
		})
	})

	player.RepeatModeChanged.OnLazy(func(rm player.RepeatMode) bool {
		switch rm {
		case player.RepeatModeNone:
			repeatClassState.SetValue("")
			repeatIconState.SetValue(repeatListIconName)
		case player.RepeatModeQueue:
			repeatClassState.SetValue("accent")
			repeatIconState.SetValue(repeatListIconName)
		case player.RepeatModeTrack:
			repeatClassState.SetValue("accent")
			repeatIconState.SetValue(repeatTrackIconName)
		}
		return signals.Continue
	})

	player.ShuffleStateChanged.OnLazy(func(b bool) bool {
		if b {
			shuffleClassState.SetValue("accent")
		} else {
			shuffleClassState.SetValue("")
		}
		return signals.Continue
	})

	return HStack(
		controlButton.
			TooltipText(gettext.Get("Toggle Shuffle")).
			IconName("playlist-shuffle-symbolic").
			ActionName("win.player.shuffle").
			BindCSSClass(shuffleClassState),
		Spacer().VExpand(false),
		controlButton.
			TooltipText(gettext.Get("Previous")).
			IconName("skip-backward-large-symbolic").
			ActionName("win.player.previous"),
		Spacer().VExpand(false),
		Button().
			TooltipText(gettext.Get("Play / Pause")).
			ActionName("win.player.play-pause").
			BindChild(actualPlayPauseChildState).
			BindSensitive(isControllableState).
			WithCSSClass("suggested-action").CornerRadius(21).
			HPadding(32).VPadding(9),
		Spacer().VExpand(false),
		controlButton.
			TooltipText(gettext.Get("Next")).
			IconName("skip-forward-large-symbolic").
			ActionName("win.player.next"),
		Spacer().VExpand(false),
		controlButton.
			TooltipText(gettext.Get("Toggle Repeat")).
			BindIconName(repeatIconState).BindCSSClass(repeatClassState).
			ActionName("win.player.repeat"),
	).HExpand(true)
}
