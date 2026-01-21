package routes

import (
	"context"
	"fmt"
	"log/slog"
	"unsafe"

	"codeberg.org/dergs/tonearm/internal/notifications"
	"codeberg.org/dergs/tonearm/internal/player"
	"codeberg.org/dergs/tonearm/internal/resources"
	"codeberg.org/dergs/tonearm/internal/router"
	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/internal/ui/components/tracklist"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/pagination"
	"codeberg.org/dergs/tonearm/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gdk"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var canPlayPlaylistState = state.NewStateful(false)

func init() {
	player.PlaybackStateChanged.On(func(ps *player.PlaybackState) bool {
		canPlayPlaylistState.SetValue(!ps.Loading)
		return signals.Continue
	})

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
		if artwork.Attributes.IsPicture() {
			coverUrl = artwork.Attributes.Files.AtLeast(160).Href
			break
		}
	}

	list := tracklist.NewTrackList(
		tracklist.GroupedColumn(2, gtk.AlignStartValue, tracklist.CoverColumn, tracklist.TitleAlbumColumn),
		tracklist.ArtistsColumn,
		tracklist.ExpandCustomButtonColumn(1, func(trackId string, position, _ int) {
			go player.PlayPlaylist(playlistUUID, false, position)
		}),
		tracklist.GroupedColumn(1, gtk.AlignEndValue, tracklist.DurationColumn, tracklist.ControlsColumn),
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
							IconName("playlist-shuffle-symbolic").
							MinWidth(81).
							CornerRadius(21).
							Padding(9).
							VAlign(gtk.AlignCenterValue).
							ConnectClicked(func(b gtk.Button) {
								go player.PlayPlaylist(playlistUUID, true, 0)
							}).
							BindSensitive(canPlayPlaylistState),
						Button().
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
							VAlign(gtk.AlignCenterValue).
							ConnectClicked(func(b gtk.Button) {
								go player.PlayPlaylist(playlistUUID, false, 0)
							}).
							BindSensitive(canPlayPlaylistState),
					).
						Spacing(5).
						HAlign(gtk.AlignEndValue),
					HStack(
						Button().
							IconName("heart-outline-thick-symbolic").
							WithCSSClass("transparent"),
						Button().
							IconName("share-alt-symbolic").
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
									var list *tracklist.TrackList[*openapi.Track]
									list = (*tracklist.TrackList[*openapi.Track])(unsafe.Pointer(u))
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
