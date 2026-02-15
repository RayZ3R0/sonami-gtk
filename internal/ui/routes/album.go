package routes

import (
	"strings"

	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/player"
	"codeberg.org/dergs/tonearm/internal/resources"
	"codeberg.org/dergs/tonearm/internal/router"
	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
	"codeberg.org/dergs/tonearm/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var canPlayAlbumState = state.NewStateful(false)

func init() {
	router.Register("album/:id", Album)
	player.PlaybackStateChanged.On(func(ps *player.PlaybackState) bool {
		canPlayAlbumState.SetValue(!ps.Loading)
		return signals.Continue
	})
}

func Album(albumId string) *router.Response {
	service, err := injector.Inject[tonearm.Service]()
	if err != nil {
		return router.FromError(gettext.Get("Album"), err)
	}

	album, err := service.GetAlbum(albumId, tonearm.AlbumHintArtists, tonearm.AlbumHintCover, tonearm.AlbumHintTracks)
	if err != nil {
		return router.FromError(gettext.Get("Album"), err)
	}

	coverUrl, err := album.Cover(160)
	if err != nil {
		return router.FromError(gettext.Get("Album"), err)
	}

	artistPaginator, err := album.Artists()
	if err != nil {
		return router.FromError(gettext.Get("Album"), err)
	}

	artists, err := artistPaginator.GetAll()
	if err != nil {
		return router.FromError(gettext.Get("Album"), err)
	}

	artistNames := []string{}
	for _, artist := range artists {
		artistNames = append(artistNames, artist.Name())
	}

	// paginator := pagination.NewPaginator(tidal.OpenAPI.V2.Albums.Items, albumId, func(items *openapi.Response[[]openapi.Relationship]) []openapi.Track {
	// 	return items.Included.Tracks(items.Data...)
	// }, "items", "items.artists", "items.albums.coverArt")

	// album, err := tidal.OpenAPI.V2.Albums.Album(context.Background(), albumId, "coverArt", "artists", "coverArt")
	// if err != nil {
	// 	return router.FromError(gettext.Get("Album"), err)
	// }

	// artists := []string{}
	// for _, artist := range album.Included.PlainArtists(album.Data.Relationships.Artists.Data...) {
	// 	artists = append(artists, artist.Attributes.Name)
	// }

	// coverUrl := ""
	// for _, artwork := range album.Included.PlainArtworks(album.Data.Relationships.CoverArt.Data...) {
	// 	if artwork.Attributes.IsPicture() {
	// 		coverUrl = artwork.Attributes.Files.AtLeast(160).Href
	// 		break
	// 	}
	// }

	// page, err := pages.NewPaginatedTracklistPage(
	// 	paginator,
	// 	func() *tracklist.TrackList[*openapi.Track] {
	// 		return tracklist.NewTrackList(
	// 			tracklist.GroupedColumn(2, gtk.AlignStartValue, tracklist.PositionColumn, tracklist.TitleColumn),
	// 			tracklist.ArtistsColumn,
	// 			tracklist.ExpandCustomButtonColumn(1, func(trackId string, position, _ int) {
	// 				go player.PlayAlbum(albumId, false, position)
	// 			}),
	// 			tracklist.GroupedColumn(1, gtk.AlignEndValue, tracklist.DurationColumn, tracklist.ControlsColumn),
	// 		)
	// 	}, func(tl *tracklist.TrackList[*openapi.Track]) schwifty.BaseWidgetable {
	// 		return tl.HMargin(30).VAlign(gtk.AlignStartValue)
	// 	},
	// )

	return &router.Response{
		PageTitle: album.Title(),
		Error:     err,
		View: VStack(
			HStack(
				AspectFrame(
					Image().
						PixelSize(146).
						FromPaintable(resources.MissingAlbum()).
						ConnectConstruct(func(i *gtk.Image) {
							if coverUrl != "" {
								injector.MustInject[*imgutil.ImgUtil]().LoadIntoImage(coverUrl, i)
							}
						}),
				).CornerRadius(10).Overflow(gtk.OverflowHiddenValue),
				VStack(
					Label(album.Title()).
						WithCSSClass("title-2").
						HAlign(gtk.AlignStartValue),
					Label(strings.Join(artistNames, ", ")).
						WithCSSClass("heading").WithCSSClass("dimmed").
						PaddingTop(10).
						HAlign(gtk.AlignStartValue),
					// Label(album.Data.Attributes.ReleaseDate.Format("2006")).
					// 	WithCSSClass("heading").WithCSSClass("dimmed").
					// 	HAlign(gtk.AlignStartValue),
					// Label(gettext.GetN("%d Track (%s)", "%d Tracks (%s)", album.Data.Attributes.NumberOfItems, album.Data.Attributes.NumberOfItems, tidalapi.FormatDuration(album.Data.Attributes.Duration.Duration))).
					// 	WithCSSClass("heading").WithCSSClass("dimmed").
					// 	HAlign(gtk.AlignStartValue).
					// 	MarginTop(10),
				).MarginStart(20).VAlign(gtk.AlignCenterValue),
				Spacer().VExpand(false),
				HStack(
					Button().
						TooltipText(gettext.Get("Shuffle Album")).
						IconName("playlist-shuffle-symbolic").
						MinWidth(81).
						CornerRadius(21).
						Padding(9).
						ConnectClicked(func(b gtk.Button) {
							go player.PlayAlbum(albumId, true, 0)
						}).
						BindSensitive(canPlayAlbumState),
					Button().
						TooltipText(gettext.Get("Play Album")).
						IconName("play-symbolic").
						MinWidth(81).
						CornerRadius(21).
						Padding(9).
						CSS(`
							button {
								background-color: var(--accent-bg-color);
							}

							button:hover {
								background-color: var(--accent-color);
							}
						`).
						ConnectClicked(func(b gtk.Button) {
							go player.PlayAlbum(albumId, false, 0)
						}).
						BindSensitive(canPlayAlbumState),
				).
					VAlign(gtk.AlignCenterValue).
					Spacing(5),
			).HMargin(40),
			// page.VExpand(true).MarginTop(20),
		).VMargin(20),
	}
}
