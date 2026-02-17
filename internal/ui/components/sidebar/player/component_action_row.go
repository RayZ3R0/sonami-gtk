package player

import (
	"context"
	"slices"
	"strconv"

	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/notifications"
	"codeberg.org/dergs/tonearm/internal/player"
	"codeberg.org/dergs/tonearm/internal/router"
	"codeberg.org/dergs/tonearm/internal/secrets"
	"codeberg.org/dergs/tonearm/internal/signals"
	appState "codeberg.org/dergs/tonearm/internal/state"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/schwifty/tracking"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gdk"
	"github.com/jwijenbergh/puregotk/v4/gio"
	"github.com/jwijenbergh/puregotk/v4/glib"
	"github.com/jwijenbergh/puregotk/v4/gobject"
	"github.com/jwijenbergh/puregotk/v4/gtk"
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
	player.TrackChanged.OnLazy(func(t tonearm.Track) bool {
		isTrackLoadedState.SetValue(t != nil)

		if t != nil {
			artists := t.Artists()
			if len(artists) > 1 {
				menu := gio.NewMenu()
				defer menu.Unref()
				for _, artist := range artists {
					menu.AppendItem(gio.NewMenuItem(artist.Title(), "win.route.artist::"+artist.ID()))
				}
				artistButtonState.SetValue(artistButtonMultiple.MenuModel(&menu.MenuModel))
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

	player.TrackChanged.On(func(t tonearm.Track) bool {
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
			weakRef := tracking.NewWeakRef(&b.Object)

			isLoading.On(func(loading bool) bool {
				schwifty.OnMainThreadOncePure(func() {
					weakRef.Use(func(obj *gobject.Object) {
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
					weakRef.Use(func(obj *gobject.Object) {
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

func actionRow() schwifty.Box {
	return HStack(
		MenuButton().
			TooltipText(gettext.Get("Volume")).
			Popover(controlsVolumeSlider()).
			IconName("speakers-symbolic").
			WithCSSClass("flat"),
		favouriteButton(),
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
				notifications.OnToast.Notify(gettext.Get("Copied track URL to clipboard."))
			}),
	).HAlign(gtk.AlignCenterValue).Spacing(15).CSS("box { margin-top: -8px; margin-bottom: -8px; }")
}
