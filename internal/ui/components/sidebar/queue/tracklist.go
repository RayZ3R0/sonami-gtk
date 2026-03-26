package queue

import (
	"slices"

	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/internal/player"
	"github.com/RayZ3R0/sonami-gtk/internal/player/queue"
	"github.com/RayZ3R0/sonami-gtk/internal/signals"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components/tracklist"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/state"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
)

var (
	baseQueueState = state.NewStateful([]sonami.Track{})
	userQueueState = state.NewStateful([]sonami.Track{})
)

func init() {
	player.BaseQueue.Entries().On(func(tracks []sonami.Track) bool {
		schwifty.OnMainThreadOnce(func(u uintptr) {
			baseQueueState.SetValue(tracks)
			filledState := queueFilledState.Value()
			filledState.Base = len(tracks) > 0
			queueFilledState.SetValue(filledState)
		}, 0)
		return signals.Continue
	})
	player.UserQueue.Entries().On(func(tracks []sonami.Track) bool {
		schwifty.OnMainThreadOnce(func(u uintptr) {
			userQueueState.SetValue(tracks)
			filledState := queueFilledState.Value()
			filledState.User = len(tracks) > 0
			queueFilledState.SetValue(filledState)
		}, 0)
		return signals.Continue
	})

}

func makeQueueTracklist(q queue.Queue, queueState *state.State[[]sonami.Track]) *tracklist.TrackList {
	trackList := tracklist.NewTrackList(
		tracklist.CoverColumn, tracklist.TitleAlbumColumn,
		tracklist.CustomWidgetButtonColumn(func(_ string, position int, _ int32) *gtk.Widget {
			return Button().
				TooltipText(gettext.Get("Remove Track from Queue")).
				IconName("user-trash-symbolic").
				WithCSSClass("flat").
				MarginEnd(10).
				ConnectClicked(func(b gtk.Button) {
					go q.RemoveAt(position)
				}).
				ToGTK()
		}),
	)

	trackList.SetClickHandler(func(track sonami.Track, position int) {
		go player.SkipThroughQueue(q, position)
	})

	trackList.BindTracks(queueState)

	trackList.SetReorderCallback(func(sourceIndex, targetIndex int, track sonami.Track) {
		q.Entries().Notify(func(oldValue []sonami.Track) []sonami.Track {
			q := slices.Clone(oldValue)
			q = append(q[:sourceIndex], q[sourceIndex+1:]...)
			q = append(q[:targetIndex], append([]sonami.Track{track}, q[targetIndex:]...)...)
			return q
		})
	})

	return trackList
}
