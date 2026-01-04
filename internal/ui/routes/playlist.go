package routes

import (
	"context"
	"fmt"
	"log/slog"
	"unsafe"

	"codeberg.org/dergs/tidalwave/internal/notifications"
	"codeberg.org/dergs/tidalwave/internal/player"
	"codeberg.org/dergs/tidalwave/internal/resources"
	"codeberg.org/dergs/tidalwave/internal/router"
	"codeberg.org/dergs/tidalwave/internal/ui/components/tracklist"
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/pagination"
	"codeberg.org/dergs/tidalwave/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gdk"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func init() {
	router.Register("playlist/:id", Playlist)
}

func Playlist(playlistUUID string) *router.Response {
	tidal := injector.MustInject[*tidalapi.TidalAPI]()

	playlist, err := tidal.OpenAPI.V2.Playlists.Playlist(context.Background(), playlistUUID, "coverArt", "ownerProfiles")
	if err != nil {
		return router.FromError("Playlist", err)
	}

	paginator := pagination.NewPaginator(tidal.OpenAPI.V2.Playlists, playlistUUID, func(r *openapi.Response[[]openapi.Relationship]) []openapi.Track {
		return r.Included.Tracks(r.Data...)
	}, "items", "items.artists", "items.albums.coverArt")
	items, err := paginator.GetFirstPage()
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
		tracklist.CustomButtonColumn(func(trackId string) {
			go player.PlayPlaylist(playlistUUID, false, trackId)
		}),
		tracklist.ControlsColumn,
	)

	for _, track := range items {
		list.AddTrack(&track)
	}

	var playlistMetadata schwifty.Label
	if playlist.Data.Attributes.PlaylistType != openapi.PlaylistTypeMix {
		playlistMetadata = Label(fmt.Sprintf("%d Tracks (%s)", playlist.Data.Attributes.NumberOfItems, tidalapi.FormatDuration(playlist.Data.Attributes.Duration.Duration)))
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
				VStack(
					HStack(
						Button().
							IconName("media-playlist-shuffle-symbolic").
							MinWidth(81).
							CornerRadius(21).
							Padding(9).
							VAlign(gtk.AlignCenterValue).
							ConnectClicked(func(b gtk.Button) {
								go player.PlayPlaylist(playlistUUID, true, "")
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
								go player.PlayPlaylist(playlistUUID, false, "")
							}),
					).
						Spacing(5).
						HAlign(gtk.AlignEndValue),
					HStack(
						Button().
							IconName("heart-outline-thick-symbolic").
							WithCSSClass("transparent"),
						Button().
							IconName("folder-publicshare-symbolic").
							WithCSSClass("transparent").
							ConnectClicked(func(gtk.Button) {
								id := playlist.Data.ID
								if id == "" {
									notifications.OnToast.Notify("No playlist could be shared.")
									return
								}

								display := gdk.DisplayGetDefault()
								defer display.Unref()
								clipboard := display.GetClipboard()
								defer clipboard.Unref()

								clipboard.SetText(fmt.Sprintf("https://tidal.com/playlist/%s", id))
								notifications.OnToast.Notify("Copied playlist URL to clipboard.")
							}),
					).
						Spacing(10).
						HAlign(gtk.AlignEndValue),
				).
					Spacing(20).
					VAlign(gtk.AlignCenterValue),
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
