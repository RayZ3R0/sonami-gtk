package favouritebutton

import (
	"log/slog"
	"slices"

	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/internal/state"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	gtkbindings "codeberg.org/dergs/tonearm/pkg/schwifty/bindings/gtk"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/schwifty/tracking"
	"github.com/jwijenbergh/puregotk/v4/gobject"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var logger = slog.With("module", "ui/components", "component", "FavouriteButton")

func FavouriteButton(favouriteCache state.FavouriteCache, resourceID string) gtkbindings.Button {
	isFavourited := signals.NewStatefulSignal(false)

	return Button().
		TooltipText(gettext.Get("Add to Collection")).
		IconName("heart-outline-thick-symbolic").
		WithCSSClass("flat").
		ConnectConstruct(func(b *gtk.Button) {
			isFavourited.Notify(func(oldValue bool) bool {
				items, err := favouriteCache.Get()
				if err != nil {
					logger.Error("Failed to get favourites", "error", err)
					return false
				}

				return slices.Contains(*items, resourceID)
			})

			weakRef := tracking.NewWeakRef(&b.Object)
			isFavourited.On(func(value bool) bool {
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
			isFavourited.Notify(func(oldValue bool) bool {
				if oldValue {
					err := favouriteCache.Remove(resourceID)
					if err != nil {
						logger.Error("error while removing item from favourites", "error", err)
						return oldValue
					}
				} else {
					err := favouriteCache.Add(resourceID)
					if err != nil {
						logger.Error("error while adding item to favourites", "error", err)
						return oldValue
					}
				}

				return !oldValue
			})
		})
}
