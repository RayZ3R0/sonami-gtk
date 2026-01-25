package routes

import (
	"context"
	"fmt"

	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/notifications"
	"codeberg.org/dergs/tonearm/internal/resources"
	"codeberg.org/dergs/tonearm/internal/router"
	"codeberg.org/dergs/tonearm/internal/ui/components"
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
						FontSize(18).
						FontWeight(700).
						HAlign(gtk.AlignStartValue),
					Label(gettext.GetN("%d Fan", "%d Fans", artistPage.Header.FollowersAmount, artistPage.Header.FollowersAmount)).
						FontSize(14).
						FontWeight(600).
						HAlign(gtk.AlignStartValue),
					Label(artistPage.Header.Biography.Text).
						FontSize(16).
						FontWeight(500).
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
							IconName("playlist-shuffle-symbolic").
							MinWidth(81).
							CornerRadius(21).
							Padding(9).
							VAlign(gtk.AlignCenterValue).
							Sensitive(false).
							ConnectClicked(func(b gtk.Button) {
								// go player.PlayPlaylist(playlistUUID, true, "")
							}),
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
							Sensitive(false).
							ConnectClicked(func(b gtk.Button) {
								// go player.PlayPlaylist(playlistUUID, false, "")
							}),
					).
						Spacing(5).
						HAlign(gtk.AlignEndValue),
					HStack(
						Button().
							IconName("heart-outline-thick-symbolic").
							WithCSSClass("transparent").
							Sensitive(false),
						Button().
							IconName("share-alt-symbolic").
							WithCSSClass("transparent").
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
