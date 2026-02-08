package routes

import (
	"context"
	"fmt"
	"log/slog"
	"slices"

	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/notifications"
	"codeberg.org/dergs/tonearm/internal/player"
	"codeberg.org/dergs/tonearm/internal/resources"
	"codeberg.org/dergs/tonearm/internal/router"
	"codeberg.org/dergs/tonearm/internal/secrets"
	"codeberg.org/dergs/tonearm/internal/signals"
	appState "codeberg.org/dergs/tonearm/internal/state"
	"codeberg.org/dergs/tonearm/internal/ui/components"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/schwifty/tracking"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	"codeberg.org/dergs/tonearm/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gdk"
	"github.com/jwijenbergh/puregotk/v4/gobject"
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

	isArtistFavourited := signals.NewStatefulSignal(false)

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
						HAlign(gtk.AlignStartValue),
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
								go player.PlayArtistTopSongs(artistId, true, 0)
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
								go player.PlayArtistTopSongs(artistId, false, 0)
							}),
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
								favList := favLists.Artist
								if err != nil {
									albumLogger.Error("Failed to load favourites", err)
									b.SetIconName("heart-outline-thick-symbolic")
									b.RemoveCssClass("accent")

									return
								}

								isArtistFavourited.Notify(func(oldValue bool) bool {
									return slices.Contains(favList, artistId)
								})

								weakRef := tracking.NewWeakRef(&b.Object)
								isArtistFavourited.On(func(value bool) bool {
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

								isArtistFavourited.Notify(func(oldValue bool) bool {
									if oldValue {
										err := tidal.V1.Favourites.RemoveArtist(context.Background(), secrets.UserID(), artistId)
										if err != nil {
											artistLogger.Error("error while removing artist from favourites", "error", err)
											return oldValue
										}
									} else {
										err := tidal.V1.Favourites.AddArtist(context.Background(), secrets.UserID(), artistId)
										if err != nil {
											artistLogger.Error("error while adding artist to favourites", "error", err)
											return oldValue
										}
									}

									return !oldValue
								})
							}), Button().
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
