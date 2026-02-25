package tracklist_header

import (
	"codeberg.org/dergs/tonearm/internal/player"
	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"github.com/jwijenbergh/puregotk/v4/gtk"
	"github.com/jwijenbergh/puregotk/v4/pango"
)

var playerControlsAvailableState = state.NewStateful(false)

func init() {
	player.PlaybackStateChanged.On(func(ps *player.PlaybackState) bool {
		playerControlsAvailableState.SetValue(!ps.Loading)
		return signals.Continue
	})
}

func template(coverUrl string, title string, subtitle string, description string, controls schwifty.Box, secondaryControls schwifty.Box) schwifty.Box {
	return HStack(
		componentCover(coverUrl),
		VStack(
			// Title + Subtitle / Main Playback Controls
			HStack(
				componentHeading(title, subtitle),
				controls,
			),
			// Description / Secondary Controls
			HStack(
				Label(description).Wrap(true).Lines(3).Ellipsis(pango.EllipsizeEndValue).WithCSSClass("dimmed").HAlign(gtk.AlignStartValue).HExpand(true),
				secondaryControls.MarginStart(10),
			).MarginTop(20),
		).MarginStart(20).VAlign(gtk.AlignCenterValue),
	).WithCSSClass("tracklist_header").HExpand(true)
}
