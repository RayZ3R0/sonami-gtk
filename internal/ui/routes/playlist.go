package routes

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"

	"codeberg.org/dergs/tidalwave/internal/player"
	"codeberg.org/dergs/tidalwave/internal/resources"
	"codeberg.org/dergs/tidalwave/internal/router"
	"codeberg.org/dergs/tidalwave/internal/ui/components/tracklist"
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	"codeberg.org/dergs/tidalwave/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func init() {
	router.Register("playlist", Playlist)
}

func Playlist(params router.Params) *router.Response {
	playlistUUID, ok := params["id"].(string)
	if !ok {
		return router.FromError("Playlist", errors.New("invalid playlist ID"))
	}

	tidal := injector.MustInject[*tidalapi.TidalAPI]()

	playlist, err := tidal.OpenAPI.V2.Playlists.Playlist(context.Background(), playlistUUID, "coverArt", "ownerProfiles")
	if err != nil {
		return router.FromError("Playlist", err)
	}

	// TODO: Handle pagination with scroll events
	items, err := tidal.OpenAPI.V2.Playlists.Items(context.Background(), playlistUUID, "", "items", "items.artists", "items.albums", "items.albums.coverArt")
	if err != nil {
		return router.FromError("Playlist", err)
	}

	creator := "TIDAL"
	for _, artist := range playlist.Included.PlainArtists(playlist.Data.Relationships.OwnerProfiles.Data...) {
		creator = artist.Attributes.Name
		break
	}

	coverUrl := ""
	for _, artwork := range playlist.Included.PlainArtworks(playlist.Data.Relationships.CoverArt.Data...) {
		coverUrl = artwork.Attributes.Files.AtLeast(320).Href
	}

	list := tracklist.NewTrackList(
		"",
		tracklist.CoverColumn,
		tracklist.TitleAlbumColumn,
		tracklist.ArtistsColumn,
		tracklist.DurationColumn,
		tracklist.ButtonColumn,
		tracklist.ControlsColumn,
	)

	for _, track := range items.Included.Tracks(items.Data...) {
		list.AddTrack(&track)
	}

	var playlistMetadata schwifty.Label
	if playlist.Data.Attributes.PlaylistType != openapi.PlaylistTypeMix {
		playlistMetadata = Label(fmt.Sprintf("%d Tracks (%s)", playlist.Data.Attributes.NumberOfItems, tidalapi.FormatDuration(int(playlist.Data.Attributes.Duration.Seconds()))))
	} else {
		playlistMetadata = Label("Personal Mix")
	}

	return &router.Response{
		PageTitle: playlist.Data.Attributes.Name,
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
					Label(playlist.Data.Attributes.Name).
						FontSize(18).
						FontWeight(700).
						HAlign(gtk.AlignStartValue),
					Label(creator).
						FontSize(16).
						FontWeight(500).
						HAlign(gtk.AlignStartValue),
					Label(playlist.Data.Attributes.CreatedAt.Format("2006")).
						FontSize(16).
						FontWeight(500).
						HAlign(gtk.AlignStartValue),
					playlistMetadata.
						FontSize(14).
						FontWeight(600).
						HAlign(gtk.AlignStartValue).
						MarginTop(10),
				).MarginStart(20).VAlign(gtk.AlignCenterValue),
				Spacer().
					VExpand(false),
				HStack(
					Button().
						IconName("media-playlist-shuffle-symbolic").
						MinWidth(81).
						CornerRadius(21).
						Padding(9).
						VAlign(gtk.AlignCenterValue).
						ConnectClicked(func(b gtk.Button) {
							i := rand.IntN(len(items.Data))
							player.Play(items.Included.Tracks(items.Data[i])[0].Data.ID)
						}),
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
						VAlign(gtk.AlignCenterValue).
						ConnectClicked(func(b gtk.Button) {
							player.Play(items.Included.Tracks(items.Data[0])[0].Data.ID)
						}),
				).
					Spacing(5),
			),
			ScrolledWindow().
				Child(list).
				Policy(gtk.PolicyNeverValue, gtk.PolicyAutomaticValue).
				VExpand(true).
				MarginTop(20),
		).HMargin(40).VMargin(20),
	}
}
