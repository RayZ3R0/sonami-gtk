package queue

import (
	"codeberg.org/dergs/tidalwave/internal/player"
	"codeberg.org/dergs/tidalwave/internal/signals"
	"codeberg.org/dergs/tidalwave/internal/ui/components/tracklist"
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var baseQueueState = state.NewStateful([]*openapi.Track{})
var userQueueState = state.NewStateful([]*openapi.Track{})

func init() {
	player.OnBaseQueueChanged.On(func(tracks []*openapi.Track) bool {
		schwifty.OnMainThreadOnce(func(u uintptr) {
			baseQueueState.SetValue(tracks)
		}, 0)
		return signals.Continue
	})
	player.OnUserQueueChanged.On(func(tracks []*openapi.Track) bool {
		schwifty.OnMainThreadOnce(func(u uintptr) {
			userQueueState.SetValue(tracks)
		}, 0)
		return signals.Continue
	})
}

func NewQueue() schwifty.ScrolledWindow {
	trackList := tracklist.NewTrackList("", tracklist.CoverColumn, tracklist.TitleAlbumColumn, tracklist.DurationColumn)
	trackList.BindTracks(userQueueState)

	trackListBase := tracklist.NewTrackList("", tracklist.CoverColumn, tracklist.TitleAlbumColumn, tracklist.DurationColumn)
	trackListBase.BindTracks(baseQueueState)
	return ScrolledWindow().
		HMargin(10).
		Child(
			VStack(
				trackList.Background("alpha(var(--view-bg-color), 0.9)").CornerRadius(10),
				trackListBase,
			).Spacing(10).VAlign(gtk.AlignStartValue),
		)
}
