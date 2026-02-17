package my_collection

import (
	"log/slog"

	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/notifications"
	"codeberg.org/dergs/tonearm/internal/player"
	"codeberg.org/dergs/tonearm/internal/resources"
	"codeberg.org/dergs/tonearm/internal/router"
	"codeberg.org/dergs/tonearm/internal/secrets"
	"codeberg.org/dergs/tonearm/internal/services/tidal/openapi"
	"codeberg.org/dergs/tonearm/internal/ui/components"
	"codeberg.org/dergs/tonearm/internal/ui/components/tracklist"
	"codeberg.org/dergs/tonearm/internal/ui/pages"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	modelopenapi "codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
	"codeberg.org/dergs/tonearm/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gtk"
	"github.com/jwijenbergh/puregotk/v4/pango"
)

func Tracks() *router.Response {
	logger := slog.With("module", "ui").WithGroup("ui").With("route", "my_collection").WithGroup("my_collection")
	tidal := injector.MustInject[*tidalapi.TidalAPI]()
	userId := secrets.UserID()
	if userId == "" {
		return &router.Response{
			PageTitle: gettext.Get("My Collection"),
			View:      components.AuthRequired(gettext.Get("Please sign in to view your collection")),
		}
	}

	paginator := openapi.NewPaginator(tidal.OpenAPI.V2.UserCollections.Tracks, userId, func(r *modelopenapi.Response[[]modelopenapi.Relationship]) []tonearm.Track {
		results := r.Included.Tracks(r.Data...)
		tracks := make([]tonearm.Track, len(results))
		for i, track := range results {
			tracks[i] = openapi.NewTrack(track)
		}
		return tracks
	}, "tracks.artists", "tracks.albums.coverArt")

	page, err := pages.NewPaginatedTracklistPage(paginator, func() *tracklist.TrackList {
		return tracklist.NewTrackList(
			tracklist.GroupedColumn(2, gtk.AlignStartValue, tracklist.CoverColumn, tracklist.TitleAlbumColumn),
			tracklist.ArtistsColumn,
			tracklist.ExpandButtonColumn(1),
			tracklist.GroupedColumn(1, gtk.AlignEndValue, tracklist.DurationColumn, tracklist.ControlsColumn),
		)
	}, func(tl *tracklist.TrackList) schwifty.BaseWidgetable {
		return tl.HMargin(40).VAlign(gtk.AlignStartValue)
	})

	return &router.Response{
		PageTitle: gettext.Get("My Tracks"),
		Error:     err,
		View: VStack(
			components.MainContent(
				HStack(
					AspectFrame(
						Image().
							PixelSize(146).
							FromPaintable(resources.MissingAlbum()).
							ConnectConstruct(func(i *gtk.Image) {
								injector.MustInject[*imgutil.ImgUtil]().LoadIntoImage(new(openapi.MyTracksInfo).Cover(146), i)
							}),
					).
						CornerRadius(10).
						Overflow(gtk.OverflowHiddenValue),
					Label(gettext.Get("My Tracks")).WithCSSClass("title-1").Ellipsis(pango.EllipsizeEndValue),
					Spacer().VExpand(false).MinWidth(20),
					HStack(
						Button().
							TooltipText(gettext.Get("Shuffle Album")).
							IconName("playlist-shuffle-symbolic").
							WithCSSClass("pill").
							ConnectClicked(func(b gtk.Button) {
								go func() {
									tracks, err := paginator.GetAll()
									if err != nil {
										notifications.OnToast.Notify(gettext.Get("An error occurred while shuffling the tracks"))
										logger.Error("An error occurred while fetching the tracks", "error", err.Error())
										return
									}

									if err := player.PlayTracklist(new(openapi.MyTracksInfo), tracks, true, 0); err != nil {
										notifications.OnToast.Notify(gettext.Get("An error occurred while shuffling the tracks"))
										logger.Error("An error occurred while playing the tracks", "error", err.Error())
									}
								}()
							}),
						Button().
							TooltipText(gettext.Get("Play Album")).
							IconName("play-symbolic").
							WithCSSClass("pill").
							WithCSSClass("suggested-action").
							ConnectClicked(func(b gtk.Button) {
								go func() {
									tracks, err := paginator.GetAll()
									if err != nil {
										notifications.OnToast.Notify(gettext.Get("An error occurred while playing the tracks"))
										logger.Error("An error occurred while fetching the tracks", "error", err.Error())
										return
									}

									if err := player.PlayTracklist(new(openapi.MyTracksInfo), tracks, false, 0); err != nil {
										notifications.OnToast.Notify(gettext.Get("An error occurred while playing the tracks"))
										logger.Error("An error occurred while playing the tracks", "error", err.Error())
									}
								}()
							}),
					).
						Spacing(12).
						VAlign(gtk.AlignCenterValue),
				).Spacing(20).HMargin(40),
			),
			page.
				VExpand(true),
		).
			VMargin(20).
			Spacing(20),
	}
}
