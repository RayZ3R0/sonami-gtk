package queue

import (
	"log/slog"
	"slices"
	"strings"
	"time"

	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/player"
	"codeberg.org/dergs/tonearm/internal/resources"
	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/internal/ui/components/tracklist"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
	"codeberg.org/dergs/tonearm/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gtk"
	"github.com/jwijenbergh/puregotk/v4/pango"
)

const (
	queueScrollMargin = 75
)

var baseQueueState = state.NewStateful([]tonearm.Track{})
var userQueueState = state.NewStateful([]tonearm.Track{})

var (
	coverState    = state.NewStateful[schwifty.Paintable](resources.MissingAlbum())
	trackTitle    = state.NewStateful[string]("")
	trackArtists  = state.NewStateful[string]("")
	playPauseIcon = state.NewStateful("play-symbolic")
)

var miniPlayerCanControl = state.NewStateful(false)

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
	player.PlaybackStateChanged.On(func(ps *player.PlaybackState) bool {
		miniPlayerCanControl.SetValue(!ps.Loading)
		return signals.Continue
	})
}

func NewQueue() schwifty.Box {
	trackList := tracklist.NewTrackList(
		tracklist.GroupedColumn(3, gtk.AlignStartValue, tracklist.CoverColumn, tracklist.TitleAlbumColumn),
		tracklist.ExpandCustomButtonColumn(1, func(_ string, position, _ int) {
			go player.SkipThroughQueue(player.UserQueue, position)
		}),
		tracklist.GroupedColumn(1, gtk.AlignCenterValue,
			tracklist.CustomWidgetButtonColumn(func(_ string, position, _ int) *gtk.Widget {
				return Button().
					TooltipText(gettext.Get("Remove Track from Queue")).
					IconName("user-trash-symbolic").
					WithCSSClass("flat").
					ConnectClicked(func(b gtk.Button) {
						go player.UserQueue.RemoveAt(position)
					}).
					ToGTK()
			}),
		),
	)
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
		tracklist.GroupedColumn(3, gtk.AlignStartValue, tracklist.CoverColumn, tracklist.TitleAlbumColumn),
		tracklist.ExpandCustomButtonColumn(1, func(_ string, position, _ int) {
			go player.SkipThroughQueue(player.BaseQueue, position)
		}),
		tracklist.GroupedColumn(1, gtk.AlignCenterValue,
			tracklist.CustomWidgetButtonColumn(func(_ string, position, _ int) *gtk.Widget {
				return Button().
					TooltipText(gettext.Get("Remove Track from Queue")).
					IconName("user-trash-symbolic").
					WithCSSClass("flat").
					ConnectClicked(func(b gtk.Button) {
						go player.BaseQueue.RemoveAt(position)
					}).
					ToGTK()
			}),
		),
	)
	trackListBase.BindTracks(baseQueueState)
	trackListBase.SetReorderCallback(func(sourceIndex, targetIndex int, track tonearm.Track) {
		player.BaseQueue.Entries().Notify(func(oldValue []tonearm.Track) []tonearm.Track {
			q := slices.Clone(oldValue)
			q = append(q[:sourceIndex], q[sourceIndex+1:]...)
			q = append(q[:targetIndex], append([]tonearm.Track{track}, q[targetIndex:]...)...)
			return q
		})
	})

	player.TrackChanged.On(func(trackInfo tonearm.Track) bool {
		if trackInfo != nil {
			coverUrl := trackInfo.Cover(80)
			if coverUrl == "" {
				slog.Error("Failed to load cover URL")
				return signals.Continue
			}

			if texture, err := injector.MustInject[*imgutil.ImgUtil]().Load(coverUrl); err == nil {
				schwifty.OnMainThreadOncePure(func() {
					coverState.SetValue(texture)
					texture.Unref()
				})
			}

			schwifty.OnMainThreadOncePure(func() {
				trackTitle.SetValue(tonearm.FormatTitle(trackInfo))
				trackArtists.SetValue(strings.Join(trackInfo.Artists().Names(), ", "))
			})
		}

		return signals.Continue
	})

	player.PlaybackStateChanged.On(func(state *player.PlaybackState) bool {
		schwifty.OnMainThreadOncePure(func() {
			switch state.Status {
			case player.PlaybackStatusPaused, player.PlaybackStatusStopped:
				playPauseIcon.SetValue("play-symbolic")
			case player.PlaybackStatusPlaying:
				playPauseIcon.SetValue("pause-symbolic")
			}
		})
		return signals.Continue
	})

	var miniPlayerLoadingIconSub *signals.Subscription

	return VStack(
		HStack(
			AspectFrame(
				Image().
					PixelSize(54).
					BindPaintable(coverState),
			).
				Overflow(gtk.OverflowHiddenValue).
				Background("alpha(var(--view-fg-color), 0.1)").
				CornerRadius(6),
			VStack(
				Label("").
					BindText(trackTitle).
					FontWeight(600).
					Ellipsis(pango.EllipsizeEndValue).
					HAlign(gtk.AlignStartValue),
				Label("").
					BindText(trackArtists).
					Ellipsis(pango.EllipsizeEndValue).
					HAlign(gtk.AlignStartValue),
			).
				VAlign(gtk.AlignCenterValue),
			Spacer().VExpand(false),
			HStack(
				Button().
					TooltipText(gettext.Get("Play / Pause")).
					WithCSSClass("flat").
					BindIconName(playPauseIcon).
					ConnectClicked(func(b gtk.Button) {
						player.PlayPause()
					}).
					BindSensitive(miniPlayerCanControl).
					ConnectConstruct(func(b *gtk.Button) {
						ptr := b.GoPointer()
						miniPlayerLoadingIconSub = player.PlaybackStateChanged.On(func(ps *player.PlaybackState) bool {
							if ps.Loading {
								schwifty.OnMainThreadOncePure(func() {
									b := gtk.ButtonNewFromInternalPtr(ptr)
									child := Spinner().ToGTK()
									b.SetChild(child)
								})
							}
							return signals.Continue
						})
					}).
					ConnectDestroy(func(w gtk.Widget) {
						player.PlaybackStateChanged.Unsubscribe(miniPlayerLoadingIconSub)
					}),
				Button().
					TooltipText(gettext.Get("Next")).
					WithCSSClass("flat").
					IconName("skip-forward-large-symbolic").
					ActionName("win.player.next").
					BindSensitive(miniPlayerCanControl),
			).
				Spacing(7).
				HAlign(gtk.AlignEndValue).
				VAlign(gtk.AlignCenterValue),
		).
			Spacing(16).
			Padding(12).
			MarginBottom(12).
			MarginTop(12).
			HMargin(16).
			Background("alpha(var(--view-fg-color), 0.1)").
			CornerRadius(12),
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
