package routes

import (
	"context"
	"fmt"

	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/notifications"
	"codeberg.org/dergs/tonearm/internal/player"
	"codeberg.org/dergs/tonearm/internal/resources"
	"codeberg.org/dergs/tonearm/internal/router"
	"codeberg.org/dergs/tonearm/internal/signals"
	appState "codeberg.org/dergs/tonearm/internal/state"
	favouritebutton "codeberg.org/dergs/tonearm/internal/ui/components/favourite_button"
	"codeberg.org/dergs/tonearm/internal/ui/components/tracklist"
	"codeberg.org/dergs/tonearm/internal/ui/pages"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/pagination"
	"codeberg.org/dergs/tonearm/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gdk"
	"github.com/jwijenbergh/puregotk/v4/gio"
	"github.com/jwijenbergh/puregotk/v4/glib"
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
	var appCache appState.FavouriteCache
	if playlist.Data.Attributes.PlaylistType != openapi.PlaylistTypeMix {
		playlistMetadata = Label(gettext.GetN("%d Track (%s)", "%d Tracks (%s)", playlist.Data.Attributes.NumberOfItems, playlist.Data.Attributes.NumberOfItems, tidalapi.FormatCustomDuration(playlist.Data.Attributes.Duration)))
		appCache = appState.PlaylistsCache
	} else {
		playlistMetadata = Label(gettext.Get("Personal Mix"))
		appCache = appState.MixesCache
	}

	page, err := pages.NewPaginatedTracklistPage(
		paginator,
		func() *tracklist.TrackList[*openapi.Track] {
			return tracklist.NewTrackList(
				tracklist.GroupedColumn(2, gtk.AlignStartValue, tracklist.CoverColumn, tracklist.TitleAlbumColumn),
				tracklist.ArtistsColumn,
				tracklist.ExpandCustomButtonColumn(1, func(trackId string, position, _ int) {
					go func() {
						if err := player.PlayPlaylist(playlistUUID, false, position); err != nil {
							notifications.OnToast.Notify(gettext.Get("An error occurred while playing the track"))
							albumLogger.Error("An error occurred while playing the playlist", "error", err.Error())
						}
					}()
				}),
				tracklist.GroupedColumn(1, gtk.AlignEndValue, tracklist.DurationColumn, tracklist.ControlsColumn),
			)
		}, func(tl *tracklist.TrackList[*openapi.Track]) schwifty.BaseWidgetable {
			return tl.HMargin(30).VAlign(gtk.AlignStartValue)
		},
	)

	playControlsMenu := gio.NewMenu()
	queueAllItem := gio.NewMenuItem("Add playlist to queue", "win.player.queue")
	queueAllItem.SetActionAndTargetValue("win.player.queue", glib.NewVariantString(fmt.Sprintf("playlist/%s", playlistUUID)))
	playControlsMenu.AppendItem(queueAllItem)
	playControlsPopover := gtk.NewPopoverMenuFromModel(&playControlsMenu.MenuModel)

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
							WithCSSClass("pill").
							VAlign(gtk.AlignCenterValue).
							ConnectClicked(func(b gtk.Button) {
								go func() {
									if err := player.PlayPlaylist(playlistUUID, true, 0); err != nil {
										notifications.OnToast.Notify(gettext.Get("An error occurred while playing the playlist"))
										albumLogger.Error("An error occurred while playing the playlist", "error", err.Error())
									}
								}()
							}).
							BindSensitive(canPlayPlaylistState),
						Button().
							TooltipText(gettext.Get("Play Playlist")).
							IconName("play-symbolic").
							WithCSSClass("pill").
							WithCSSClass("suggested-action").
							VAlign(gtk.AlignCenterValue).
							ConnectClicked(func(b gtk.Button) {
								go func() {
									if err := player.PlayPlaylist(playlistUUID, false, 0); err != nil {
										notifications.OnToast.Notify(gettext.Get("An error occurred while playing the playlist"))
										albumLogger.Error("An error occurred while playing the playlist", "error", err.Error())
									}
								}()
							}).
							BindSensitive(canPlayPlaylistState),
						MenuButton().
							TooltipText(gettext.Get("More…")).
							WithCSSClass("circular").
							WithCSSClass("flat").
							VAlign(gtk.AlignCenterValue).
							IconName("view-more-symbolic").
							Popover(playControlsPopover),
					).
						Spacing(12).
						HAlign(gtk.AlignEndValue),
					HStack(
						favouritebutton.FavouriteButton(appCache, playlistUUID),
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
