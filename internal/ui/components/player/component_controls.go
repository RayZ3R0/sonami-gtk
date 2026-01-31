package player

import (
	"codeberg.org/dergs/tonearm/internal/g"
	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/player"
	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

const (
	repeatListIcon  = "playlist-repeat-symbolic"
	repeatTrackIcon = "playlist-repeat-song-symbolic"
)

var pauseIcon = g.Lazy(func() any {
	return Image().FromIconName("pause-symbolic")()
})
var playIcon = g.Lazy(func() any {
	return Image().FromIconName("play-symbolic")()
})
var playSpinner = g.Lazy(func() any {
	return Spinner().SizeRequest(16, 16)()
})

var (
	playPauseChildState = state.NewStateful[any](nil)

	repeatClassState = state.NewStateful("")
	repeatIconState  = state.NewStateful(repeatListIcon)

	shuffleClassState = state.NewStateful("")
)

var controlButton = Button().
	BindSensitive(isControllableState).
	MinHeight(34).MinWidth(34).
	WithCSSClass("transparent").
	VAlign(gtk.AlignCenterValue)

func init() {
	wasLoading := false
	player.PlaybackStateChanged.On(func(state *player.PlaybackState) bool {
		schwifty.OnMainThreadOncePure(func() {
			if state.Loading && !wasLoading {
				playPauseChildState.SetValue(playSpinner())
				wasLoading = true
			} else {
				wasLoading = false
				switch state.Status {
				case player.PlaybackStatusPlaying:
					playPauseChildState.SetValue(pauseIcon())
				case player.PlaybackStatusPaused, player.PlaybackStatusStopped:
					playPauseChildState.SetValue(playIcon())
				}
			}
		})

		return signals.Continue
	})

	player.RepeatModeChanged.OnLazy(func(rm player.RepeatMode) bool {
		switch rm {
		case player.RepeatModeNone:
			repeatClassState.SetValue("")
			repeatIconState.SetValue(repeatListIcon)
		case player.RepeatModeQueue:
			repeatClassState.SetValue("accent")
			repeatIconState.SetValue(repeatListIcon)
		case player.RepeatModeTrack:
			repeatClassState.SetValue("accent")
			repeatIconState.SetValue(repeatTrackIcon)
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
}

func controls() schwifty.Box {
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
			BindChild(playPauseChildState).
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
