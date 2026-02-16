package routes

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

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
	"github.com/jwijenbergh/puregotk/v4/pango"
)

var albumLogger = slog.With("module", "ui/routes", "route", "album")
var canPlayAlbumState = state.NewStateful(false)

func init() {
	router.Register("album/:id", Album)
	player.PlaybackStateChanged.On(func(ps *player.PlaybackState) bool {
		canPlayAlbumState.SetValue(!ps.Loading)
		return signals.Continue
	})
}

func Album(albumId string) *router.Response {
	tidal := injector.MustInject[*tidalapi.TidalAPI]()

	paginator := pagination.NewPaginator(tidal.OpenAPI.V2.Albums.Items, albumId, func(items *openapi.Response[[]openapi.Relationship]) []openapi.Track {
		return items.Included.Tracks(items.Data...)
	}, "items", "items.artists", "items.albums.coverArt")

	album, err := tidal.OpenAPI.V2.Albums.Album(context.Background(), albumId, "coverArt", "artists", "coverArt")
	if err != nil {
		return router.FromError(gettext.Get("Album"), err)
	}

	artists := []string{}
	for _, artist := range album.Included.PlainArtists(album.Data.Relationships.Artists.Data...) {
		artists = append(artists, artist.Attributes.Name)
	}

	coverUrl := ""
	for _, artwork := range album.Included.PlainArtworks(album.Data.Relationships.CoverArt.Data...) {
		if artwork.Attributes.IsPicture() {
			coverUrl = artwork.Attributes.Files.AtLeast(160).Href
			break
		}
	}

	page, err := pages.NewPaginatedTracklistPage(
		paginator,
		func() *tracklist.TrackList[*openapi.Track] {
			return tracklist.NewTrackList(
				tracklist.GroupedColumn(2, gtk.AlignStartValue, tracklist.PositionColumn, tracklist.TitleColumn),
				tracklist.ArtistsColumn,
				tracklist.ExpandCustomButtonColumn(1, func(trackId string, position, _ int) {
					go func() {
						if err := player.PlayAlbum(albumId, false, position); err != nil {
							notifications.OnToast.Notify(gettext.Get("An error occurred while playing the track"))
							albumLogger.Error("An error occurred while playing the album", "error", err.Error())
						}
					}()
				}),
				tracklist.GroupedColumn(1, gtk.AlignEndValue, tracklist.DurationColumn, tracklist.ControlsColumn),
			)
		}, func(tl *tracklist.TrackList[*openapi.Track]) schwifty.BaseWidgetable {
			return tl.HMargin(30).VAlign(gtk.AlignStartValue)
		},
	)

	if err != nil {
		return router.FromError(album.Data.Attributes.Title, err)
	}

	playControlsMenu := gio.NewMenu()
	queueAllItem := gio.NewMenuItem("Add album to queue", "win.player.queue")
	queueAllItem.SetActionAndTargetValue("win.player.queue", glib.NewVariantString(fmt.Sprintf("album/%s", albumId)))
	playControlsMenu.AppendItem(queueAllItem)
	playControlsPopover := gtk.NewPopoverMenuFromModel(&playControlsMenu.MenuModel)

	return &router.Response{
		PageTitle: album.Data.Attributes.Title,
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
					Label(album.Data.Attributes.Title).
						WithCSSClass("title-2").
						HAlign(gtk.AlignStartValue).
						Ellipsis(pango.EllipsizeEndValue),
					Label(strings.Join(artists, ", ")).
						Ellipsis(pango.EllipsizeEndValue).
						WithCSSClass("heading").WithCSSClass("dimmed").
						PaddingTop(10).
						HAlign(gtk.AlignStartValue),
					Label(album.Data.Attributes.ReleaseDate.Format("2006")).
						WithCSSClass("heading").WithCSSClass("dimmed").
						HAlign(gtk.AlignStartValue),
					Label(gettext.GetN("%d Track (%s)", "%d Tracks (%s)", album.Data.Attributes.NumberOfItems, album.Data.Attributes.NumberOfItems, tidalapi.FormatDuration(album.Data.Attributes.Duration.Duration))).
						WithCSSClass("heading").WithCSSClass("dimmed").
						HAlign(gtk.AlignStartValue).
						MarginTop(10),
				).
					MarginStart(20).
					VAlign(gtk.AlignCenterValue),
				Spacer().VExpand(false),
				VStack(
					HStack(
						Button().
							TooltipText(gettext.Get("Shuffle Album")).
							IconName("playlist-shuffle-symbolic").
							WithCSSClass("pill").
							ConnectClicked(func(b gtk.Button) {
								go func() {
									if err := player.PlayAlbum(albumId, true, 0); err != nil {
										notifications.OnToast.Notify(gettext.Get("An error occurred while playing the album"))
										albumLogger.Error("An error occurred while playing the album", "error", err.Error())
									}
								}()
							}).
							BindSensitive(canPlayAlbumState),
						Button().
							TooltipText(gettext.Get("Play Album")).
							IconName("play-symbolic").
							WithCSSClass("pill").
							WithCSSClass("suggested-action").
							ConnectClicked(func(b gtk.Button) {
								go func() {
									if err := player.PlayAlbum(albumId, false, 0); err != nil {
										notifications.OnToast.Notify(gettext.Get("An error occurred while playing the album"))
										albumLogger.Error("An error occurred while playing the album", "error", err.Error())
									}
								}()
							}).
							BindSensitive(canPlayAlbumState),
						MenuButton().
							TooltipText(gettext.Get("More…")).
							Popover(playControlsPopover).
							WithCSSClass("flat").
							WithCSSClass("circular").
							IconName("view-more-symbolic"),
					).
						VAlign(gtk.AlignCenterValue).
						Spacing(12),
					HStack(
						favouritebutton.FavouriteButton(appState.AlbumsCache, albumId),
						Button().
							TooltipText(gettext.Get("Copy Album URL")).
							IconName("share-alt-symbolic").
							WithCSSClass("flat").
							ConnectClicked(func(gtk.Button) {
								id := album.Data.ID
								if id == "" {
									notifications.OnToast.Notify(gettext.Get("No album could be shared."))
									return
								}

								display := gdk.DisplayGetDefault()
								defer display.Unref()
								clipboard := display.GetClipboard()
								defer clipboard.Unref()

								clipboard.SetText(fmt.Sprintf("https://tidal.com/album/%s?u", id))
								notifications.OnToast.Notify(gettext.Get("Copied album URL to clipboard."))
							}),
					).
						Spacing(10).
						HAlign(gtk.AlignEndValue),
				).
					MarginStart(20).
					HAlign(gtk.AlignEndValue).
					VAlign(gtk.AlignCenterValue),
			).
				HMargin(40),
			page.VExpand(true).MarginTop(20),
		).VMargin(20),
	}
}
