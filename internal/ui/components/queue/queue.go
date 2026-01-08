package queue

import (
	"log/slog"

	"codeberg.org/dergs/tidalwave/internal/player"
	"codeberg.org/dergs/tidalwave/internal/resources"
	"codeberg.org/dergs/tidalwave/internal/signals"
	"codeberg.org/dergs/tidalwave/internal/ui/components/tracklist"
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	"codeberg.org/dergs/tidalwave/pkg/utils/imgutil"
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
	playPauseIcon = state.NewStateful("media-playback-start-symbolic")
)

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
}

func NewQueue() schwifty.Box {
	trackList := tracklist.NewTrackList("Queued tracks", tracklist.CoverColumn, tracklist.TitleAlbumColumn, tracklist.DurationColumn)
	trackList.BindTracks(userQueueState)

	trackListBase := tracklist.NewTrackList("Coming up", tracklist.CoverColumn, tracklist.TitleAlbumColumn, tracklist.DurationColumn)
	trackListBase.BindTracks(baseQueueState)

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
			case player.PlaybackStatusBuffering, player.PlaybackStatusPaused:
				playPauseIcon.SetValue("media-playback-start-symbolic")
			case player.PlaybackStatusPlaying:
				playPauseIcon.SetValue("media-playback-pause-symbolic")
			}
		})
		return signals.Continue
	})

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
					}),
				Button().
					WithCSSClass("transparent").
					IconName("media-skip-forward-symbolic").
					ActionName("win.player.next"),
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
