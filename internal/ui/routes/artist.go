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

	return &router.Response{
		PageTitle: gettext.Get("Artist"),
		View: VStack(
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
							MinWidth(81).
							CornerRadius(21).
							Padding(9).
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
								go func() {
									if err := player.PlayArtistTopSongs(artistId, false, 0); err != nil {
										notifications.OnToast.Notify(gettext.Get("An error occurred while playing the top tracks"))
										albumLogger.Error("An error occurred while playing the top tracks", "error", err.Error())
									}
								}()
							}),
					).
						Spacing(5).
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
					VAlign(gtk.AlignCenterValue),
			).HMargin(40),
			ScrolledWindow().
				Child(body).
				Policy(gtk.PolicyNeverValue, gtk.PolicyAutomaticValue).
				VExpand(true).
				MarginTop(20),
		).VMargin(20),
	}
}
