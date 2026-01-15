package queue

import (
	"log/slog"
	"slices"

	"codeberg.org/dergs/tonearm/internal/player"
	"codeberg.org/dergs/tonearm/internal/resources"
	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/internal/ui/components/tracklist"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"
	"codeberg.org/dergs/tonearm/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gtk"
	"github.com/jwijenbergh/puregotk/v4/pango"
)

var baseQueueState = state.NewStateful([]*openapi.Track{})
var userQueueState = state.NewStateful([]*openapi.Track{})

var (
	coverState    = state.NewStateful[schwifty.Paintable](resources.MissingAlbum())
	trackTitle    = state.NewStateful[string]("")
	trackArtists  = state.NewStateful[string]("")
	playPauseIcon = state.NewStateful("play-symbolic")
)

var miniPlayerCanControl = state.NewStateful(false)

var log = slog.With("module", "queue")

func init() {
	player.BaseQueue.UpcomingEntries.On(func(tracks []*openapi.Track) bool {
		schwifty.OnMainThreadOnce(func(u uintptr) {
			baseQueueState.SetValue(tracks)
		}, 0)
		return signals.Continue
	})
	player.UserQueue.UpcomingEntries.On(func(tracks []*openapi.Track) bool {
		schwifty.OnMainThreadOnce(func(u uintptr) {
			userQueueState.SetValue(tracks)
		}, 0)
		return signals.Continue
	})
	player.ControllableStateChanged.On(func(cs player.ControllableState) bool {
		miniPlayerCanControl.SetValue(cs.CanControl())
		return signals.Continue
	})
}

func NewQueue() schwifty.Box {
	trackList := tracklist.NewTrackList(
		tracklist.GroupedColumn(3, gtk.AlignStartValue, tracklist.CoverColumn, tracklist.TitleAlbumColumn),
		tracklist.ExpandCustomButtonColumn(1, func(_ string, position, _ int) {
			go player.SkipThoughQueue(player.UserQueue, position)
		}),
		tracklist.GroupedColumn(1, gtk.AlignCenterValue,
			tracklist.CustomWidgetButtonColumn(func(_ string, position, _ int) *gtk.Widget {
				return Button().
					IconName("user-trash-symbolic").
					WithCSSClass("transparent").
					ConnectClicked(func(b gtk.Button) {
						go player.UserQueue.Remove(position)
					}).
					ToGTK()
			}),
		),
	)
	trackList.BindTracks(userQueueState)
	trackList.SetReorderCallback(func(sourceIndex, targetIndex int, track *openapi.Track) {
		player.UserQueue.UpcomingEntries.Notify(func(oldValue []*openapi.Track) []*openapi.Track {
			q := slices.Clone(oldValue)
			q = append(q[:sourceIndex], q[sourceIndex+1:]...)
			q = append(q[:targetIndex], append([]*openapi.Track{track}, q[targetIndex:]...)...)
			return q
		})
	})

	trackListBase := tracklist.NewTrackList(
		tracklist.GroupedColumn(3, gtk.AlignStartValue, tracklist.CoverColumn, tracklist.TitleAlbumColumn),
		tracklist.ExpandCustomButtonColumn(1, func(_ string, position, _ int) {
			go player.SkipThoughQueue(player.BaseQueue, position)
		}),
		tracklist.GroupedColumn(1, gtk.AlignCenterValue,
			tracklist.CustomWidgetButtonColumn(func(_ string, position, _ int) *gtk.Widget {
				return Button().
					IconName("user-trash-symbolic").
					WithCSSClass("transparent").
					ConnectClicked(func(b gtk.Button) {
						go player.BaseQueue.Remove(position)
					}).
					ToGTK()
			}),
		),
	)
	trackListBase.BindTracks(baseQueueState)
	trackListBase.SetReorderCallback(func(sourceIndex, targetIndex int, track *openapi.Track) {
		player.BaseQueue.UpcomingEntries.Notify(func(oldValue []*openapi.Track) []*openapi.Track {
			q := slices.Clone(oldValue)
			q = append(q[:sourceIndex], q[sourceIndex+1:]...)
			q = append(q[:targetIndex], append([]*openapi.Track{track}, q[targetIndex:]...)...)
			return q
		})
	})

	player.TrackChanged.On(func(trackInfo *player.Track) bool {
		if trackInfo != nil {
			if trackInfo.CoverURL != "" {
				if texture, err := injector.MustInject[*imgutil.ImgUtil]().Load(trackInfo.CoverURL); err == nil {
					schwifty.OnMainThreadOncePure(func() {
						coverState.SetValue(texture)
						texture.Unref()
					})
				}
			}

			schwifty.OnMainThreadOncePure(func() {
				trackTitle.SetValue(trackInfo.Title)
				trackArtists.SetValue(trackInfo.ArtistNames())
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
					WithCSSClass("transparent").
					BindIconName(playPauseIcon).
					ConnectClicked(func(b gtk.Button) {
						player.PlayPause()
					}).
					BindSensitive(miniPlayerCanControl).
					ConnectConstruct(func(b *gtk.Button) {
						ptr := b.GoPointer()
						miniPlayerLoadingIconSub = player.ControllableStateChanged.On(func(cs player.ControllableState) bool {
							if !cs.PlayerReady {
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
						player.ControllableStateChanged.Unsubscribe(miniPlayerLoadingIconSub)
					}),
				Button().
					WithCSSClass("transparent").
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
			VPadding(10).
			VExpand(true).
			Child(
				VStack(
					trackList.Background("alpha(var(--view-bg-color), 0.9)").CornerRadius(10),
					trackListBase,
				).Spacing(10).VAlign(gtk.AlignStartValue),
			),
	)
}
