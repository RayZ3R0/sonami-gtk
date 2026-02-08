package favouritebutton

import (
	"context"
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
	"github.com/jwijenbergh/puregotk/v4/gobject"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var logger = slog.With("module", "ui/components", "component", "FavouriteButton")

func FavouriteButton(favList []string, resourceID string, apiEndpoint interface {
	Add(context.Context, string, string) error
	Remove(context.Context, string, string) error
}) gtkbindings.Button {
	isFavourited := signals.NewStatefulSignal(false)

	return Button().
		TooltipText(gettext.Get("Add to Collection")).
		IconName("heart-outline-thick-symbolic").
		WithCSSClass("flat").
		ConnectConstruct(func(b *gtk.Button) {
			favLists, err := state.Favourites()
			favList := favLists.Album
			if err != nil {
				logger.Error("Failed to load favourites", err)
				b.SetIconName("heart-outline-thick-symbolic")
				b.RemoveCssClass("accent")

				return
			}

			isFavourited.Notify(func(oldValue bool) bool {
				return slices.Contains(favList, resourceID)
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
					err := apiEndpoint.Remove(context.Background(), secrets.UserID(), resourceID)
					if err != nil {
						logger.Error("error while removing album from favourites", "error", err)
						return oldValue
					}
				} else {
					err := apiEndpoint.Add(context.Background(), secrets.UserID(), resourceID)
					if err != nil {
						logger.Error("error while adding album to favourites", "error", err)
						return oldValue
					}
				}

				return !oldValue
			})
		})
}
