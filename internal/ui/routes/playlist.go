package routes

import (
	"context"
	"fmt"
	"slices"

	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/notifications"
	"codeberg.org/dergs/tonearm/internal/player"
	"codeberg.org/dergs/tonearm/internal/resources"
	"codeberg.org/dergs/tonearm/internal/router"
	"codeberg.org/dergs/tonearm/internal/secrets"
	"codeberg.org/dergs/tonearm/internal/signals"
	appState "codeberg.org/dergs/tonearm/internal/state"
	"codeberg.org/dergs/tonearm/internal/ui/components/tracklist"
	"codeberg.org/dergs/tonearm/internal/ui/pages"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/schwifty/tracking"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/pagination"
	"codeberg.org/dergs/tonearm/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gdk"
	"github.com/jwijenbergh/puregotk/v4/gobject"
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
		return router.FromError(gettext.Get("Playlist"), err)
	}

	paginator := pagination.NewPaginator(tidal.OpenAPI.V2.Playlists.Items, playlistUUID, func(r *openapi.Response[[]openapi.Relationship]) []openapi.Track {
		return r.Included.Tracks(r.Data...)
	}, "items", "items.artists", "items.albums.coverArt")

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

	var playlistMetadata schwifty.Label
	if playlist.Data.Attributes.PlaylistType != openapi.PlaylistTypeMix {
		playlistMetadata = Label(gettext.GetN("%d Track (%s)", "%d Tracks (%s)", playlist.Data.Attributes.NumberOfItems, playlist.Data.Attributes.NumberOfItems, tidalapi.FormatCustomDuration(playlist.Data.Attributes.Duration)))
	} else {
		playlistMetadata = Label(gettext.Get("Personal Mix"))
	}

	page, err := pages.NewPaginatedTracklistPage(
		paginator,
		func() *tracklist.TrackList[*openapi.Track] {
			return tracklist.NewTrackList(
				tracklist.GroupedColumn(2, gtk.AlignStartValue, tracklist.CoverColumn, tracklist.TitleAlbumColumn),
				tracklist.ArtistsColumn,
				tracklist.ExpandCustomButtonColumn(1, func(trackId string, position, _ int) {
					go player.PlayPlaylist(playlistUUID, false, position)
				}),
				tracklist.GroupedColumn(1, gtk.AlignEndValue, tracklist.DurationColumn, tracklist.ControlsColumn),
			)
		}, func(tl *tracklist.TrackList[*openapi.Track]) schwifty.BaseWidgetable {
			return tl.HMargin(30).VAlign(gtk.AlignStartValue)
		},
	)

	isPlaylistFavourited := signals.NewStatefulSignal(false)

	return &router.Response{
		PageTitle: playlist.Data.Attributes.Name,
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
					Label(playlist.Data.Attributes.Name).
						WithCSSClass("title-2").
						HAlign(gtk.AlignStartValue),
					Label(creator).
						WithCSSClass("heading").WithCSSClass("dimmed").
						PaddingTop(10).
						HAlign(gtk.AlignStartValue),
					Label(playlist.Data.Attributes.CreatedAt.Format("2006")).
						WithCSSClass("heading").WithCSSClass("dimmed").
						HAlign(gtk.AlignStartValue),
					playlistMetadata.
						WithCSSClass("heading").WithCSSClass("dimmed").
						HAlign(gtk.AlignStartValue).
						MarginTop(10),
				).MarginStart(20).VAlign(gtk.AlignCenterValue),
				Spacer().
					VExpand(false),
				VStack(
					HStack(
						Button().
							TooltipText(gettext.Get("Shuffle Playlist")).
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
							TooltipText(gettext.Get("Play Playlist")).
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
							TooltipText(gettext.Get("Add to Collection")).
							IconName("heart-outline-thick-symbolic").
							WithCSSClass("flat").
							ConnectConstruct(func(b *gtk.Button) {
								favLists, err := appState.Favourites()
								favList := favLists.Playlist
								if err != nil {
									albumLogger.Error("Failed to load favourites", err)
									b.SetIconName("heart-outline-thick-symbolic")
									b.RemoveCssClass("accent")

									return
								}

								isPlaylistFavourited.Notify(func(oldValue bool) bool {
									return slices.Contains(favList, playlistUUID)
								})

								weakRef := tracking.NewWeakRef(&b.Object)
								isPlaylistFavourited.On(func(value bool) bool {
									schwifty.OnMainThreadOncePure(func() {
										weakRef.Use(func(obj *gobject.Object) {
											b := gtk.ButtonNewFromInternalPtr(obj.Ptr)

											if value {
												b.SetIconName("heart-filled-symbolic")
												b.AddCssClass("accent")
											} else {
												b.SetIconName("heart-outline-thick-symbolic")
												b.RemoveCssClass("accent")
											}
										})
									})

									return signals.Continue
								})
							}).
							ConnectClicked(func(b gtk.Button) {
								tidal, _ := injector.Inject[*tidalapi.TidalAPI]()

								isPlaylistFavourited.Notify(func(oldValue bool) bool {
									if oldValue {
										err := tidal.V1.Favourites.RemovePlaylist(context.Background(), secrets.UserID(), playlistUUID)
										if err != nil {
											albumLogger.Error("error while removing playlist from favourites", "error", err)
											return oldValue
										}
									} else {
										err := tidal.V1.Favourites.AddPlaylist(context.Background(), secrets.UserID(), playlistUUID)
										if err != nil {
											albumLogger.Error("error while adding playlist to favourites", "error", err)
											return oldValue
										}
									}

									return !oldValue
								})
							}),
						Button().
							TooltipText(gettext.Get("Copy Playlist URL")).
							IconName("share-alt-symbolic").
							WithCSSClass("flat").
							ConnectClicked(func(gtk.Button) {
								id := playlist.Data.ID
								if id == "" {
									notifications.OnToast.Notify(gettext.Get("No playlist could be shared."))
									return
								}

								display := gdk.DisplayGetDefault()
								defer display.Unref()
								clipboard := display.GetClipboard()
								defer clipboard.Unref()

								clipboard.SetText(fmt.Sprintf("https://tidal.com/playlist/%s", id))
								notifications.OnToast.Notify(gettext.Get("Copied playlist URL to clipboard."))
							}),
					).
						Spacing(10).
						HAlign(gtk.AlignEndValue),
				).
					Spacing(20).
					VAlign(gtk.AlignCenterValue),
			).HMargin(40),
			page.VExpand(true).MarginTop(20),
		).VMargin(20),
	}
}
