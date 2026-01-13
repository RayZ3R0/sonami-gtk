package player

import (
	"codeberg.org/dergs/tonearm/internal/player"
	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var (
	durationState       = state.NewStateful("00:00")
	positionState       = state.NewStateful("00:00")
	timelineSliderState = state.NewStateful(0.0)
)

func init() {
	player.TrackChanged.On(func(trackInfo *player.Track) bool {
		schwifty.OnMainThreadOncePure(func() {
			if trackInfo == nil {
				durationState.SetValue("00:00")
			} else {
				durationState.SetValue(tidalapi.FormatDuration(trackInfo.Duration))
			}
		})
		return signals.Continue
	})

	player.PlaybackStateChanged.On(func(state *player.PlaybackState) bool {
		schwifty.OnMainThreadOncePure(func() {
			positionState.SetValue(tidalapi.FormatDuration(state.Position))
			if state.Duration > 0 {
				durationState.SetValue(tidalapi.FormatDuration(state.Duration))
				timelineSliderState.SetValue(100.0 / float64(state.Duration) * float64(state.Position))
			} else {
				timelineSliderState.SetValue(0)
			}
		})
		return signals.Continue
	})
}

func trackTimeline() schwifty.Box {
	return VStack(
		Scale(gtk.OrientationHorizontalValue).
			BindSensitive(isControllable).
			BindValue(timelineSliderState).
			Increments(1, 1).
			HExpand(true).
			Range(0, 100).
			Value(0).
			Background("transparent").
			MarginTop(24).
			HMargin(24).
			HPadding(0).
			ConnectChangeValue(func(r gtk.Range, st gtk.ScrollType, f float64) bool {
				player.SeekToPercent(f)
				return false
			}),
		HStack(
			Label("00:00").BindText(positionState),
			Spacer().VExpand(false),
			Label("00:00").BindText(durationState),
		).HMargin(24),
	)
}
