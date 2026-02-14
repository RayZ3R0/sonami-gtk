package favouritebutton

import (
	"log/slog"
	"slices"

	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/secrets"
	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/internal/state"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	gtkbindings "codeberg.org/dergs/tonearm/pkg/schwifty/bindings/gtk"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/schwifty/tracking"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gobject"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var logger = slog.With("module", "ui/components", "component", "FavouriteButton")
var spinner = func() *gtk.Widget {
	return &adw.NewSpinner().Widget
}

func FavouriteButton(favouriteCache state.FavouriteCache, resourceID string) gtkbindings.Button {
	isFavourited := signals.NewStatefulSignal(false)
	isLoading := signals.NewStatefulSignal(false)

	return Button().
		TooltipText(gettext.Get("Add to Collection")).
		IconName("heart-outline-thick-symbolic").
		WithCSSClass("flat").
		BindSensitive(secrets.SignedInState).
		ConnectConstruct(func(b *gtk.Button) {
			weakRef := tracking.NewWeakRef(&b.Object)
			isLoading.Set(true)
			defer isLoading.Set(false)

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
							} else {
								b.SetIconName("heart-outline-thick-symbolic")
								b.RemoveCssClass("accent")
							}
						}
					})
				})

				return signals.Continue
			})

			isFavourited.Notify(func(oldValue bool) bool {
				items, err := favouriteCache.Get()
				if err != nil {
					logger.Error("Failed to get favourites", "error", err)
					return false
				}

				return slices.Contains(*items, resourceID)
			})

			isFavourited.On(func(value bool) bool {
				if isLoading.CurrentValue() {
					return signals.Continue
				}

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
			go func() {
				if isLoading.CurrentValue() {
					return
				}

				isLoading.Set(true)
				defer isLoading.Set(false)

				if isFavourited.CurrentValue() {
					err := favouriteCache.Remove(resourceID)
					if err != nil {
						logger.Error("error while removing item from favourites", "error", err)
						return
					}
				} else {
					err := favouriteCache.Add(resourceID)
					if err != nil {
						logger.Error("error while adding item to favourites", "error", err)
						return
					}
				}

				isFavourited.Notify(func(oldValue bool) bool {
					return !oldValue
				})
			}()
		})
}
