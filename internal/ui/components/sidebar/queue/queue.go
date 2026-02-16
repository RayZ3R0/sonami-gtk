package queue

import (
	"log/slog"
	"slices"
	"time"

	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/player"
	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/internal/ui/components/sidebar"
	"codeberg.org/dergs/tonearm/internal/ui/components/tracklist"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

const (
	queueScrollMargin = 75
)

var baseQueueState = state.NewStateful([]tonearm.Track{})
var userQueueState = state.NewStateful([]tonearm.Track{})

var log = slog.With("module", "queue")

func init() {
	player.BaseQueue.Entries().On(func(tracks []tonearm.Track) bool {
		schwifty.OnMainThreadOnce(func(u uintptr) {
			baseQueueState.SetValue(tracks)
		}, 0)
		return signals.Continue
	})
	player.UserQueue.Entries().On(func(tracks []tonearm.Track) bool {
		schwifty.OnMainThreadOnce(func(u uintptr) {
			userQueueState.SetValue(tracks)
		}, 0)
		return signals.Continue
	})
}

func NewQueue() schwifty.Box {
	trackLoaded := state.NewStateful(false)

	player.TrackChanged.On(func(t tonearm.Track) bool {
		trackLoaded.SetValue(t != nil)
		return signals.Continue
	})

	trackList := tracklist.NewTrackList(
		tracklist.CoverColumn, tracklist.TitleAlbumColumn,
		tracklist.CustomWidgetButtonColumn(func(_ string, position, _ int) *gtk.Widget {
			return Button().
				TooltipText(gettext.Get("Remove Track from Queue")).
				IconName("user-trash-symbolic").
				WithCSSClass("flat").
				MarginEnd(15).
				ConnectClicked(func(b gtk.Button) {
					go player.UserQueue.RemoveAt(position)
				}).
				ToGTK()
		}),
	)
	trackList.SetClickHandler(func(track tonearm.Track, position int) {
		go player.SkipThroughQueue(player.UserQueue, position)
	})
	trackList.BindTracks(userQueueState)
	trackList.SetReorderCallback(func(sourceIndex, targetIndex int, track tonearm.Track) {
		player.UserQueue.Entries().Notify(func(oldValue []tonearm.Track) []tonearm.Track {
			q := slices.Clone(oldValue)
			q = append(q[:sourceIndex], q[sourceIndex+1:]...)
			q = append(q[:targetIndex], append([]tonearm.Track{track}, q[targetIndex:]...)...)
			return q
		})
	})

	trackListBase := tracklist.NewTrackList(
		tracklist.CoverColumn, tracklist.TitleAlbumColumn,
		tracklist.CustomWidgetButtonColumn(func(_ string, position, _ int) *gtk.Widget {
			return Button().
				TooltipText(gettext.Get("Remove Track from Queue")).
				IconName("user-trash-symbolic").
				WithCSSClass("flat").
				MarginEnd(15).
				ConnectClicked(func(b gtk.Button) {
					go player.BaseQueue.RemoveAt(position)
				}).
				ToGTK()
		}),
	)
	trackListBase.SetClickHandler(func(track tonearm.Track, position int) {
		go player.SkipThroughQueue(player.BaseQueue, position)
	})
	trackListBase.BindTracks(baseQueueState)
	trackListBase.SetReorderCallback(func(sourceIndex, targetIndex int, track tonearm.Track) {
		player.BaseQueue.Entries().Notify(func(oldValue []tonearm.Track) []tonearm.Track {
			q := slices.Clone(oldValue)
			q = append(q[:sourceIndex], q[sourceIndex+1:]...)
			q = append(q[:targetIndex], append([]tonearm.Track{track}, q[targetIndex:]...)...)
			return q
		})
	})

	return VStack(
		sidebar.MiniPlayer().BindVisible(trackLoaded),
		ScrolledWindow().
			HMargin(10).
			VExpand(true).
			Child(
				VStack(
					trackList.Background("alpha(var(--view-bg-color), 0.9)").CornerRadius(10),
					trackListBase,
				).Spacing(10).VAlign(gtk.AlignStartValue),
			).
			ConnectRealize(func(sw gtk.Widget) {
				action := 0.0
				ptr := sw.GoPointer()
				controller := gtk.NewDropControllerMotion()
				controller.ConnectMotion(new(func(controller gtk.DropControllerMotion, _, y float64) {
					sw := gtk.ScrolledWindowNewFromInternalPtr(ptr)
					sw.Ref()
					defer sw.Unref()

					if y < queueScrollMargin {
						speed := (queueScrollMargin - y) / queueScrollMargin
						action = speed * -4
					} else if h := sw.GetAllocatedHeight(); (float64(h) - y) < queueScrollMargin {
						speed := (queueScrollMargin - (float64(h) - y)) / queueScrollMargin
						action = speed * 4
					} else {
						action = 0
					}

				}))
				controller.ConnectLeave(new(func(gtk.DropControllerMotion) {
					action = 0
				}))

				var fn func()

				fn = func() {
					sw := gtk.ScrolledWindowNewFromInternalPtr(ptr)
					sw.Ref()
					defer sw.Unref()

					adj := sw.GetVadjustment()
					schwifty.OnMainThreadOnce(func(u uintptr) {
						adj := gtk.AdjustmentNewFromInternalPtr(u)
						defer adj.Unref()

						adj.SetValue(adj.GetValue() + float64(action)*10)
					}, adj.GoPointer())

					time.AfterFunc(10*time.Millisecond, fn)
				}

				time.AfterFunc(10*time.Millisecond, fn)

				sw.AddController(&controller.EventController)
			}),
	)
}
