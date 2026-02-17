package queue

import (
	"log/slog"
	"slices"
	"time"

	"codeberg.org/dergs/tonearm/internal/g"
	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/player"
	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/internal/ui/components/sidebar"
	"codeberg.org/dergs/tonearm/internal/ui/components/tracklist"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/schwifty/tracking"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
	"github.com/jwijenbergh/puregotk/v4/gobject"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

const (
	queueScrollMargin = 75
)

type queueFilledType struct {
	User, Base bool
}

func (q queueFilledType) IsFilled() bool {
	return q.User || q.Base
}

var (
	baseQueueState   = state.NewStateful([]tonearm.Track{})
	userQueueState   = state.NewStateful([]tonearm.Track{})
	queueFilledState = state.NewStateful[queueFilledType](queueFilledType{})

	queueDisplay = state.NewStateful[any](StatusPage().Title(gettext.Get("No Tracks in Queue")).IconName("music-queue-empty-symbolic").VExpand(true))
)

var log = slog.With("module", "queue")

func init() {
	player.BaseQueue.Entries().On(func(tracks []tonearm.Track) bool {
		schwifty.OnMainThreadOnce(func(u uintptr) {
			baseQueueState.SetValue(tracks)
			filledState := queueFilledState.Value()
			filledState.Base = len(tracks) > 0
			queueFilledState.SetValue(filledState)
		}, 0)
		return signals.Continue
	})
	player.UserQueue.Entries().On(func(tracks []tonearm.Track) bool {
		schwifty.OnMainThreadOnce(func(u uintptr) {
			userQueueState.SetValue(tracks)
			filledState := queueFilledState.Value()
			filledState.User = len(tracks) > 0
			queueFilledState.SetValue(filledState)
		}, 0)
		return signals.Continue
	})

	queueFilledState.AddCallback(func(newValue queueFilledType) {
		if newValue.IsFilled() {
			queueDisplay.SetValue(queueList())
		} else {
			queueDisplay.SetValue(StatusPage().Title(gettext.Get("No Tracks in Queue")).IconName("music-queue-empty-symbolic").VExpand(true))
		}
	})
}

var queueList = g.Lazy(func() *gtk.ScrolledWindow {
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

	motionTicker := time.NewTicker(10 * time.Millisecond)

	return ScrolledWindow().
		HMargin(10).
		VExpand(true).
		Child(
			VStack(
				trackList.Background("alpha(var(--view-bg-color), 0.9)").CornerRadius(10),
				trackListBase,
			).Spacing(10).VAlign(gtk.AlignStartValue),
		).
		ConnectRealize(func(sw gtk.Widget) {
			ref := tracking.NewWeakRef(&sw)

			action := 0.0

			go func() {
				for {
					select {
					case _, ok := <-motionTicker.C:
						if !ok {
							return
						}

						if action == 0 {
							continue
						}

						ref.Use(func(obj *gobject.Object) {
							sw := gtk.ScrolledWindowNewFromInternalPtr(obj.Ptr)

							adj := sw.GetVadjustment()
							schwifty.OnMainThreadOnce(func(u uintptr) {
								adj := gtk.AdjustmentNewFromInternalPtr(u)
								defer adj.Unref()

								adj.SetValue(adj.GetValue() + float64(action)*10)
							}, adj.GoPointer())
						})
					}
				}
			}()

			controller := gtk.NewDropControllerMotion()
			controller.ConnectMotion(new(func(controller gtk.DropControllerMotion, _, y float64) {
				ref.Use(func(obj *gobject.Object) {
					sw := gtk.ScrolledWindowNewFromInternalPtr(obj.Ptr)
					if y < queueScrollMargin {
						speed := (queueScrollMargin - y) / queueScrollMargin
						action = speed * -4
					} else if h := sw.GetAllocatedHeight(); (float64(h) - y) < queueScrollMargin {
						speed := (queueScrollMargin - (float64(h) - y)) / queueScrollMargin
						action = speed * 4
					} else {
						action = 0
					}
				})
			}))

			controller.ConnectLeave(new(func(gtk.DropControllerMotion) {
				action = 0
			}))

			sw.AddController(&controller.EventController)
		})()
})

func NewQueue() schwifty.Box {
	trackLoaded := state.NewStateful(false)

	player.TrackChanged.On(func(t tonearm.Track) bool {
		trackLoaded.SetValue(t != nil)
		return signals.Continue
	})

	return VStack(
		sidebar.MiniPlayer().BindVisible(trackLoaded),
		Bin().BindChild(queueDisplay),
	)
}
