package favouritebutton

import (
	"log/slog"
	"slices"

	"codeberg.org/puregotk/puregotk/v4/adw"
	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/internal/secrets"
	"github.com/RayZ3R0/sonami-gtk/internal/signals"
	"github.com/RayZ3R0/sonami-gtk/internal/state"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	gtkbindings "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/bindings/gtk"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/utils/weak"
)

var logger = slog.With("module", "ui/components", "component", "FavouriteButton")
var spinner = func() *gtk.Widget {
	return &adw.NewSpinner().Widget
}

func FavouriteButton(favouriteCache state.FavouriteCache, resourceID string) gtkbindings.Button {
	isFavourited := signals.NewStatefulSignal(false)
	isLoading := signals.NewStatefulSignal(false)

	var loadingSubscription *signals.Subscription
	var favouritedSubscription *signals.Subscription
	return Button().
		IconName("heart-outline-thick-symbolic").
		WithCSSClass("flat").
		BindSensitive(secrets.SignedInState).
		ConnectRealize(func(b gtk.Widget) {
			weakRef := weak.NewWidgetRef(&b)
			isLoading.Set(true)
			defer isLoading.Set(false)

			loadingSubscription = isLoading.On(func(loading bool) bool {
				schwifty.OnMainThreadOncePure(func() {
					weakRef.Use(func(widget *gtk.Widget) {
						b := gtk.ButtonNewFromInternalPtr(widget.Ptr)

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

			isFavourited.Notify(func(oldValue bool) bool {
				items, err := favouriteCache.Get()
				if err != nil {
					logger.Error("Failed to get favourites", "error", err)
					return false
				}

				return slices.Contains(*items, resourceID)
			})

			favouritedSubscription = isFavourited.On(func(value bool) bool {
				if isLoading.CurrentValue() {
					return signals.Continue
				}

				schwifty.OnMainThreadOncePure(func() {
					weakRef.Use(func(widget *gtk.Widget) {
						b := gtk.ButtonNewFromInternalPtr(widget.Ptr)

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
		ConnectUnrealize(func(w gtk.Widget) {
			isLoading.Unsubscribe(loadingSubscription)
			isFavourited.Unsubscribe(favouritedSubscription)
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
