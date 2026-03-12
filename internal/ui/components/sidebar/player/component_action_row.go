package player

import (
	"context"
	"slices"
	"strconv"

	"codeberg.org/puregotk/puregotk/v4/adw"
	"codeberg.org/puregotk/puregotk/v4/gdk"
	"codeberg.org/puregotk/puregotk/v4/gio"
	"codeberg.org/puregotk/puregotk/v4/glib"
	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/internal/notifications"
	"github.com/RayZ3R0/sonami-gtk/internal/player"
	"github.com/RayZ3R0/sonami-gtk/internal/router"
	"github.com/RayZ3R0/sonami-gtk/internal/secrets"
	"github.com/RayZ3R0/sonami-gtk/internal/signals"
	appState "github.com/RayZ3R0/sonami-gtk/internal/state"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/state"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/utils/weak"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi"
	"github.com/infinytum/injector"
)

var (
	artistButtonSingle = Button().
				TooltipText(gettext.Get("Navigate to Artist")).
				IconName("music-artist2-symbolic").
				BindSensitive(isTrackLoadedState).
				WithCSSClass("flat")
	artistButtonMultiple = MenuButton().
				TooltipText(gettext.Get("Navigate to Artist")).
				IconName("music-artist2-symbolic").
				BindSensitive(isTrackLoadedState).
				WithCSSClass("flat")

	isTrackLoadedState = state.NewStateful(false)
	artistButtonState  = state.NewStateful[any](artistButtonSingle)
)

func init() {
	player.TrackChanged.OnLazy(func(t sonami.Track) bool {
		isTrackLoadedState.SetValue(t != nil)

		if t != nil {
			artists := t.Artists()
			if len(artists) > 1 {
				menu := gio.NewMenu()
				for _, artist := range artists {
					item := gio.NewMenuItem(artist.Title(), "win.route.artist::"+artist.ID())
					menu.AppendItem(item)
					item.Unref()
				}
				artistButtonState.SetValue(artistButtonMultiple.MenuModel(&menu.MenuModel))
				menu.Unref()
			} else {
				artistButtonState.SetValue(artistButtonSingle.ActionName("win.route.artist").ActionTargetValue(glib.NewVariantString(artists[0].ID())))
			}
		}

		return signals.Continue
	})
}

func spinner() *gtk.Widget {
	return &adw.NewSpinner().Widget
}

func favouriteButton() schwifty.Button {
	var (
		isFavourited = signals.NewStatefulSignal(false)
		isLoading    = signals.NewStatefulSignal(false)
		isSensitive  = state.NewStateful(false)
	)

	secrets.SignedInChanged.On(func(signedIn bool) bool {
		isSensitive.SetValue(player.TrackChanged.CurrentValue() != nil && signedIn)
		return signals.Continue
	})

	player.TrackChanged.On(func(t sonami.Track) bool {
		isSensitive.SetValue(t != nil && secrets.SignedInChanged.CurrentValue())
		isLoading.Set(true)
		defer isLoading.Set(false)
		if t == nil {
			return signals.Continue
		}

		favourites, err := appState.TracksCache.Get()
		if err != nil {
			logger.Error("error while fetching favourites", err)
			return signals.Continue
		}

		isFavourited.Notify(func(oldValue bool) bool {
			return slices.Contains(*favourites, t.ID())
		})

		return signals.Continue
	})

	return Button().
		IconName("heart-outline-thick-symbolic").
		WithCSSClass("flat").
		ConnectConstruct(func(b *gtk.Button) {
			weakRef := weak.NewWidgetRef(b)

			isLoading.On(func(loading bool) bool {
				schwifty.OnMainThreadOncePure(func() {
					weakRef.Use(func(obj *gtk.Widget) {
						b := gtk.ButtonNewFromInternalPtr(obj.Ptr)

						if loading {
							b.SetChild(spinner())
							b.RemoveCssClass("accent")
						} else {
							if isFavourited.CurrentValue() {
								b.SetIconName("heart-filled-symbolic")
								b.AddCssClass("accent")
								b.SetTooltipText(gettext.Get("Remove from Collection"))
							} else {
								b.SetIconName("heart-outline-thick-symbolic")
								b.RemoveCssClass("accent")
								b.SetTooltipText(gettext.Get("Add to Collection"))
							}
						}
					})
				})

				return signals.Continue
			})

			isFavourited.On(func(value bool) bool {
				schwifty.OnMainThreadOncePure(func() {
					weakRef.Use(func(obj *gtk.Widget) {
						b := gtk.ButtonNewFromInternalPtr(obj.Ptr)

						if value {
							b.SetIconName("heart-filled-symbolic")
							b.AddCssClass("accent")
							b.SetTooltipText(gettext.Get("Remove from Collection"))
						} else {
							b.SetIconName("heart-outline-thick-symbolic")
							b.RemoveCssClass("accent")
							b.SetTooltipText(gettext.Get("Add to Collection"))
						}
					})
				})

				return signals.Continue
			})
		}).
		ConnectClicked(func(b gtk.Button) {
			go func() {
				if isLoading.CurrentValue() {
					return
				}

				isLoading.Set(true)
				defer isLoading.Set(false)

				isFavourited.Notify(func(oldValue bool) bool {
					if oldValue {
						err := appState.TracksCache.Remove(player.TrackChanged.CurrentValue().ID())
						if err != nil {
							logger.Error("error while removing track from favourites", "error", err)
							return oldValue
						}
					} else {
						err := appState.TracksCache.Add(player.TrackChanged.CurrentValue().ID())
						if err != nil {
							logger.Error("error while adding track to favourites", "error", err)
							return oldValue
						}
					}

					appState.TracksCache.Bust()

					return !oldValue
				})
			}()
		}).
		BindSensitive(isSensitive)
}

func addToPlaylistButton() schwifty.Button {
	return Button().
		TooltipText(gettext.Get("Add to Playlist")).
		IconName("music-queue-symbolic").
		WithCSSClass("flat").
		BindSensitive(isTrackLoadedState).
		ConnectClicked(func(b gtk.Button) {
			track := player.TrackChanged.CurrentValue()
			if track == nil {
				return
			}

			coverURL := ""
			if album := track.Album(); album != nil {
				coverURL = album.Cover(80)
			}

			b.Widget.ActivateActionVariant("win.localplaylist.add-track", glib.NewVariantString(track.ID()+"\t"+coverURL))
		})
}

func actionRow() schwifty.Box {
	return HStack(
		MenuButton().
			TooltipText(gettext.Get("Volume")).
			Popover(controlsVolumeSlider()).
			IconName("speakers-symbolic").
			WithCSSClass("flat"),
		Button().
			TooltipText(gettext.Get("Navigate to Album")).
			IconName("cd-symbolic").
			BindSensitive(isTrackLoadedState).
			WithCSSClass("flat").
			ConnectClicked(func(b gtk.Button) {
				track := player.TrackChanged.CurrentValue()
				if track == nil {
					return
				}

				router.Navigate("album/" + track.Album().ID())
			}),
		Bin().BindChild(artistButtonState),
		Button().
			TooltipText(gettext.Get("Navigate to Track Mix")).
			IconName("compass2-symbolic").
			BindSensitive(isTrackLoadedState).
			WithCSSClass("flat").
			ConnectClicked(func(b gtk.Button) {
				track := player.TrackChanged.CurrentValue()
				if track == nil {
					return
				}

				id, _ := strconv.Atoi(track.ID())

				tidalapi := injector.MustInject[*tidalapi.TidalAPI]()
				mix, err := tidalapi.V1.Tracks.Mix(context.Background(), id)

				if err != nil {
					if err.Error() == "unexpected status code: 404" {
						notifications.OnToast.Notify(gettext.Get("No mix found for the current track"))
						return
					}

					return
				}

				router.Navigate("playlist/" + mix.ID)
			}),
		Button().
			TooltipText(gettext.Get("Copy Track URL")).
			IconName("share-alt-symbolic").
			BindSensitive(isTrackLoadedState).
			WithCSSClass("flat").
			ConnectClicked(func(gtk.Button) {
				display := gdk.DisplayGetDefault()
				defer display.Unref()
				clipboard := display.GetClipboard()
				defer clipboard.Unref()

				clipboard.SetText(player.TrackChanged.CurrentValue().URL() + "?u")
				notifications.OnToast.Notify(gettext.Get("Copied track URL to clipboard"))
			}),
		favouriteButton(),
		addToPlaylistButton(),
	).HAlign(gtk.AlignCenterValue).Spacing(15).CSS("box { margin-top: -8px; margin-bottom: -8px; }")
}
