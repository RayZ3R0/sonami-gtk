package routes

import (
	"fmt"

	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/notifications"
	"codeberg.org/dergs/tonearm/internal/player"
	"codeberg.org/dergs/tonearm/internal/resources"
	"codeberg.org/dergs/tonearm/internal/router"
	"codeberg.org/dergs/tonearm/internal/signals"
	appState "codeberg.org/dergs/tonearm/internal/state"
	"codeberg.org/dergs/tonearm/internal/ui/components"
	favouritebutton "codeberg.org/dergs/tonearm/internal/ui/components/favourite_button"
	"codeberg.org/dergs/tonearm/internal/ui/components/tracklist"
	"codeberg.org/dergs/tonearm/internal/ui/pages"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
	"codeberg.org/dergs/tonearm/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gdk"
	"github.com/jwijenbergh/puregotk/v4/gio"
	"github.com/jwijenbergh/puregotk/v4/glib"
	"github.com/jwijenbergh/puregotk/v4/gtk"
	"github.com/jwijenbergh/puregotk/v4/pango"
)

var canPlayPlaylistState = state.NewStateful(false)

func init() {
	player.PlaybackStateChanged.On(func(ps *player.PlaybackState) bool {
		canPlayPlaylistState.SetValue(!ps.Loading)
		return signals.Continue
	})

	router.Register("playlist/:id", Playlist)
}

func Playlist(playlistID string) *router.Response {
	service, err := injector.Inject[tonearm.Service]()
	if err != nil {
		return router.FromError(gettext.Get("Playlist"), err)
	}

	playlist, err := service.GetPlaylist(playlistID)
	if err != nil {
		return router.FromError(gettext.Get("Playlist"), err)
	}

	// If no creator is present, the playlist is curated by TIDAL
	creatorName := "TIDAL"
	if creator := playlist.Creator(); creator != nil {
		creatorName = creator.Title()
	}

	var playlistMetadata schwifty.Label
	var appCache appState.FavouriteCache
	if playlist.IsMix() {
		playlistMetadata = Label(gettext.Get("Personal Mix"))
		appCache = appState.MixesCache
	} else {
		playlistMetadata = Label(gettext.GetN("%d Track (%s)", "%d Tracks (%s)", playlist.Count(), playlist.Count(), tidalapi.FormatDuration(playlist.Duration())))
		appCache = appState.PlaylistsCache
	}

	trackPaginator, err := service.GetPlaylistTracks(playlistID)
	if err != nil {
		return router.FromError(gettext.Get("Playlist"), err)
	}

	page, err := pages.NewPaginatedTracklistPage(
		trackPaginator,
		func() *tracklist.TrackList {
			return tracklist.NewTrackList(
				tracklist.GroupedColumn(2, gtk.AlignStartValue, tracklist.CoverColumn, tracklist.TitleAlbumColumn),
				tracklist.ArtistsColumn,
				tracklist.ExpandCustomButtonColumn(1, func(trackId string, position, _ int) {
					go func() {
						if err := player.PlayPlaylist(playlistID, false, position); err != nil {
							notifications.OnToast.Notify(gettext.Get("An error occurred while playing the track"))
							albumLogger.Error("An error occurred while playing the playlist", "error", err.Error())
						}
					}()
				}),
				tracklist.GroupedColumn(1, gtk.AlignEndValue, tracklist.DurationColumn, tracklist.ControlsColumn),
			)
		}, func(tl *tracklist.TrackList) schwifty.BaseWidgetable {
			return tl.HMargin(30).VAlign(gtk.AlignStartValue)
		},
	)

	playControlsMenu := gio.NewMenu()
	queueAllItem := gio.NewMenuItem("Add playlist to queue", "win.player.queue")
	queueAllItem.SetActionAndTargetValue("win.player.queue", glib.NewVariantString(fmt.Sprintf("playlist/%s", playlistID)))
	playControlsMenu.AppendItem(queueAllItem)
	playControlsPopover := gtk.NewPopoverMenuFromModel(&playControlsMenu.MenuModel)

	return &router.Response{
		PageTitle: playlist.Title(),
		Error:     err,
		View: VStack(
			components.MainContent(
				HStack(
					AspectFrame(
						Image().
							PixelSize(146).
							FromPaintable(resources.MissingAlbum()).
							ConnectConstruct(func(i *gtk.Image) {
								if playlist.Cover(146) != "" {
									injector.MustInject[*imgutil.ImgUtil]().LoadIntoImage(playlist.Cover(146), i)
								}
							}),
					).CornerRadius(10).Overflow(gtk.OverflowHiddenValue),
					VStack(
						Label(playlist.Title()).
							WithCSSClass("title-2").
							Ellipsis(pango.EllipsizeEndValue).
							HAlign(gtk.AlignStartValue),
						Label(creatorName).
							Ellipsis(pango.EllipsizeEndValue).
							WithCSSClass("heading").WithCSSClass("dimmed").
							PaddingTop(10).
							HAlign(gtk.AlignStartValue),
						Label(playlist.CreatedAt().Format("2006")).
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
										if err := player.PlayPlaylist(playlistID, true, 0); err != nil {
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
										if err := player.PlayPlaylist(playlistID, false, 0); err != nil {
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
							favouritebutton.FavouriteButton(appCache, playlistID),
							Button().
								TooltipText(gettext.Get("Copy Playlist URL")).
								IconName("share-alt-symbolic").
								WithCSSClass("flat").
								ConnectClicked(func(gtk.Button) {
									display := gdk.DisplayGetDefault()
									defer display.Unref()
									clipboard := display.GetClipboard()
									defer clipboard.Unref()

									clipboard.SetText(playlist.URL())
									notifications.OnToast.Notify(gettext.Get("Copied playlist URL to clipboard."))
								}),
						).
							Spacing(10).
							HAlign(gtk.AlignEndValue),
					).
						Spacing(20).
						MarginStart(20).
						VAlign(gtk.AlignCenterValue),
				).HMargin(40),
			),
			page.VExpand(true).MarginTop(20),
		).VMargin(20),
	}
}
