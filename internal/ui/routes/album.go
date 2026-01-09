package routes

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"unsafe"

	"codeberg.org/dergs/tidalwave/internal/player"
	"codeberg.org/dergs/tidalwave/internal/resources"
	"codeberg.org/dergs/tidalwave/internal/router"
	"codeberg.org/dergs/tidalwave/internal/signals"
	"codeberg.org/dergs/tidalwave/internal/ui/components/tracklist"
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/pagination"
	"codeberg.org/dergs/tidalwave/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var canPlayAlbumState = state.NewStateful(false)

func init() {
	router.Register("album/:id", Album)
	player.ControllableStateChanged.On(func(cs player.ControllableState) bool {
		if v := cs.PlayerReady; v != canPlayAlbumState.Value() {
			canPlayAlbumState.SetValue(v)
		}
		return signals.Continue
	})
}

func Album(albumId string) *router.Response {
	tidal := injector.MustInject[*tidalapi.TidalAPI]()

	paginator := pagination.NewPaginator(tidal.OpenAPI.V2.Albums, albumId, func(items *openapi.Response[[]openapi.Relationship]) []openapi.Track {
		return items.Included.Tracks(items.Data...)
	}, "items", "items.artists", "items.albums.coverArt")

	album, err := tidal.OpenAPI.V2.Albums.Album(context.Background(), albumId, "coverArt", "artists", "coverArt")
	if err != nil {
		return router.FromError("Album", err)
	}

	items, err := paginator.GetFirstPage()
	if err != nil {
		return router.FromError("Album", err)
	}

	artists := []string{}
	for _, artist := range album.Included.PlainArtists(album.Data.Relationships.Artists.Data...) {
		artists = append(artists, artist.Attributes.Name)
	}

	coverUrl := ""
	for _, artwork := range album.Included.PlainArtworks(album.Data.Relationships.CoverArt.Data...) {
		coverUrl = artwork.Attributes.Files.AtLeast(320).Href
	}

	list := tracklist.NewTrackList(
		"",
		tracklist.PositionColumn,
		tracklist.TitleColumn,
		tracklist.ArtistsColumn,
		tracklist.DurationColumn,
		tracklist.CustomButtonColumn(func(trackId string) {
			go player.PlayAlbum(albumId, false, trackId)
		}),
		tracklist.ControlsColumn,
	)

	for _, track := range items {
		list.AddTrack(&track)
	}

	return &router.Response{
		PageTitle: album.Data.Attributes.Title,
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
					Label(album.Data.Attributes.Title).
						FontSize(18).
						FontWeight(700).
						HAlign(gtk.AlignStartValue),
					Label(strings.Join(artists, ", ")).
						FontSize(16).
						FontWeight(500).
						HAlign(gtk.AlignStartValue),
					Label(album.Data.Attributes.ReleaseDate.Format("2006")).
						FontSize(16).
						FontWeight(500).
						HAlign(gtk.AlignStartValue),
					Label(fmt.Sprintf("%d Tracks (%s)", album.Data.Attributes.NumberOfItems, tidalapi.FormatDuration(album.Data.Attributes.Duration.Duration))).
						FontSize(14).
						FontWeight(600).
						HAlign(gtk.AlignStartValue).
						MarginTop(10),
				).MarginStart(20).VAlign(gtk.AlignCenterValue),
				Spacer().VExpand(false),
				HStack(
					Button().
						IconName("media-playlist-shuffle-symbolic").
						MinWidth(81).
						CornerRadius(21).
						Padding(9).
						ConnectClicked(func(b gtk.Button) {
							go player.PlayAlbum(albumId, true, "")
						}).
						BindSensitive(canPlayAlbumState),
					Button().
						IconName("media-playback-start-symbolic").
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
							go player.PlayAlbum(albumId, false, "")
						}).
						BindSensitive(canPlayAlbumState),
				).
					VAlign(gtk.AlignCenterValue).
					Spacing(5),
			),
			ScrolledWindow().
				Child(list).
				Policy(gtk.PolicyNeverValue, gtk.PolicyAutomaticValue).
				VExpand(true).
				MarginTop(20).
				ConnectEdgeReached(func(sw gtk.ScrolledWindow, pt gtk.PositionType) {
					if pt == gtk.PosBottomValue {
						go func() {
							if !paginator.IsConsumed() {
								items, err := paginator.Next()
								if err != nil {
									return
								}

								schwifty.OnMainThreadOnce(func(u uintptr) {
									var list *tracklist.TrackList
									list = (*tracklist.TrackList)(unsafe.Pointer(u))
									for _, track := range items {
										list.AddTrack(&track)
									}
								}, uintptr(unsafe.Pointer(list)))
							} else {
								slog.Debug("No more tracks to fetch")
							}
						}()
					}
				}),
		).HMargin(40).VMargin(20),
	}
}
