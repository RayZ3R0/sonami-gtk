package player

import (
	"codeberg.org/dergs/tidalwave/internal/player"
	"codeberg.org/dergs/tidalwave/internal/signals"
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var (
	durationState       = state.NewStateful("00:00")
	positionState       = state.NewStateful("00:00")
	timelineSliderState = state.NewStateful(0.0)
)

func init() {
	player.OnTrackChanged.On(func(trackInfo player.TrackInformation) bool {
		durationState.SetValue(tidalapi.FormatDuration(int(trackInfo.Duration.Seconds())))
		return signals.Continue
	})

	player.OnStateChanged.On(func(state player.State) bool {
		positionState.SetValue(tidalapi.FormatDuration(state.Position))
		if state.Duration > 0 {
			durationState.SetValue(tidalapi.FormatDuration(int(state.Duration)))
			timelineSliderState.SetValue(100.0 / float64(state.Duration) * float64(state.Position))
		} else {
			timelineSliderState.SetValue(0)
		}

		return signals.Continue
	})
}

func trackTimeline() schwifty.Box {
	return VStack(
		Scale(gtk.OrientationHorizontalValue).
			BindValue(timelineSliderState).
			HExpand(true).
			Range(0, 100).
			Value(0).
			Background("transparent").
			MarginTop(24).
			HMargin(24).
			HPadding(0).
			ConnectChangeValue(func(r gtk.Range, st gtk.ScrollType, f float64) bool {
				player.Scrub(f)
				return false
			}),
		HStack(
			Label("00:00").BindText(positionState),
			Spacer().VExpand(false),
			Label("00:00").BindText(durationState),
		).HMargin(24),
	)
}
