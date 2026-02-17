package routes

import (
	"context"
	"fmt"
	"log/slog"

	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/notifications"
	"codeberg.org/dergs/tonearm/internal/player"
	"codeberg.org/dergs/tonearm/internal/resources"
	"codeberg.org/dergs/tonearm/internal/router"
	"codeberg.org/dergs/tonearm/internal/state"
	"codeberg.org/dergs/tonearm/internal/ui/components"
	favouritebutton "codeberg.org/dergs/tonearm/internal/ui/components/favourite_button"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	"codeberg.org/dergs/tonearm/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gdk"
	"github.com/jwijenbergh/puregotk/v4/gio"
	"github.com/jwijenbergh/puregotk/v4/glib"
	"github.com/jwijenbergh/puregotk/v4/gtk"
	"github.com/jwijenbergh/puregotk/v4/pango"
)

func init() {
	router.Register("artist/:id", Artist)
}

var artistLogger = slog.With("module", "ui/routes", "route", "artist")

func Artist(artistId string) *router.Response {
	tidal := injector.MustInject[*tidalapi.TidalAPI]()

	artistPage, err := tidal.V2.Artist.Artist(context.Background(), artistId)
	if err != nil {
		return router.FromError(gettext.Get("Artist"), err)
	}

	body := VStack().Spacing(25).VMargin(20)
	for _, item := range artistPage.Items {
		body = body.Append(components.ForPageItem(item))
	}

	playControlsMenu := gio.NewMenu()
	queueAllItem := gio.NewMenuItem(gettext.Get("Add Top Tracks to Queue"), "win.player.queue")
	queueAllItem.SetActionAndTargetValue("win.player.queue", glib.NewVariantString(fmt.Sprintf("artist/%s", artistId)))
	playControlsMenu.AppendItem(queueAllItem)
	playControlsPopover := gtk.NewPopoverMenuFromModel(&playControlsMenu.MenuModel)

	return &router.Response{
		PageTitle: gettext.Get("Artist"),
		View: VStack(
			components.MainContent(
				HStack(
					AspectFrame(
						Image().
							PixelSize(146).
							FromPaintable(resources.MissingAlbum()).
							ConnectConstruct(func(i *gtk.Image) {
								if artistPage.Item.Data.Artist.Picture != "" {
									injector.MustInject[*imgutil.ImgUtil]().LoadIntoImage(tidalapi.ImageURL(artistPage.Item.Data.Artist.Picture), i)
								}
							}),
					).CornerRadius(10).Overflow(gtk.OverflowHiddenValue),
					VStack(
						Label(artistPage.Item.Data.Artist.Name).
							WithCSSClass("title-2").
							HAlign(gtk.AlignStartValue).
							Ellipsis(pango.EllipsizeEndValue),
						Label(gettext.GetN("%d Fan", "%d Fans", artistPage.Header.FollowersAmount, artistPage.Header.FollowersAmount)).
							WithCSSClass("dimmed").
							PaddingTop(5).
							HAlign(gtk.AlignStartValue),
						Label(artistPage.Header.Biography.Text).
							WithCSSClass("dimmed").
							HAlign(gtk.AlignStartValue).
							Lines(3).
							Wrap(true).
							Ellipsis(pango.EllipsizeEndValue).
							MarginTop(10),
					).MarginStart(20).VAlign(gtk.AlignCenterValue),
					Spacer().
						VExpand(false),
					VStack(
						HStack(
							Button().
								TooltipText(gettext.Get("Shuffle Top Tracks")).
								IconName("playlist-shuffle-symbolic").
								WithCSSClass("pill").
								VAlign(gtk.AlignCenterValue).
								ConnectClicked(func(b gtk.Button) {
									go func() {
										if err := player.PlayArtistTopSongs(artistId, true, 0); err != nil {
											notifications.OnToast.Notify(gettext.Get("An error occurred while playing the top tracks"))
											albumLogger.Error("An error occurred while playing the top tracks", "error", err.Error())
										}
									}()
								}),
							Button().
								TooltipText(gettext.Get("Play Top Tracks")).
								IconName("play-symbolic").
								WithCSSClass("pill").
								WithCSSClass("suggested-action").
								VAlign(gtk.AlignCenterValue).
								ConnectClicked(func(b gtk.Button) {
									go func() {
										if err := player.PlayArtistTopSongs(artistId, false, 0); err != nil {
											notifications.OnToast.Notify(gettext.Get("An error occurred while playing the top tracks"))
											albumLogger.Error("An error occurred while playing the top tracks", "error", err.Error())
										}
									}()
								}),
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
							favouritebutton.FavouriteButton(state.ArtistsCache, artistId),
							Button().
								TooltipText(gettext.Get("Copy Artist URL")).
								IconName("share-alt-symbolic").
								WithCSSClass("flat").
								ConnectClicked(func(gtk.Button) {
									display := gdk.DisplayGetDefault()
									defer display.Unref()
									clipboard := display.GetClipboard()
									defer clipboard.Unref()

									clipboard.SetText(fmt.Sprintf("https://tidal.com/artist/%s", artistId))
									notifications.OnToast.Notify(gettext.Get("Copied artist URL to clipboard."))
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
			ScrolledWindow().
				Child(
					components.MainContent(
						body,
					),
				).
				Policy(gtk.PolicyNeverValue, gtk.PolicyAutomaticValue).
				VExpand(true).
				MarginTop(20),
		).VMargin(20),
	}
}
