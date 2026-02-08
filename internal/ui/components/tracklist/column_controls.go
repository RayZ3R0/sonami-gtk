package tracklist

import (
	"context"
	"log/slog"
	"slices"
	"strconv"

	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/secrets"
	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/internal/state"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/schwifty/tracking"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tonearm/pkg/tidalapi/models/v2"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/glib"
	"github.com/jwijenbergh/puregotk/v4/gobject"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var logger = slog.With("module", "components/tracklist")

func controlsColumn(trackId string, grid *gtk.Grid, position int, column int) int {
	isTrackFavourited := signals.NewStatefulSignal(false)

	grid.Attach(
		HStack(
			Button().
				TooltipText(gettext.Get("Add to Collection")).
				ConnectConstruct(func(b *gtk.Button) {
					favLists, err := state.Favourites()
					favList := favLists.Track
					if err != nil {
						logger.Error("Failed to load favourites", err)
						b.SetIconName("heart-outline-thick-symbolic")
						b.RemoveCssClass("accent")

						return
					}

					isTrackFavourited.Notify(func(oldValue bool) bool {
						return slices.Contains(favList, trackId)
					})

					weakRef := tracking.NewWeakRef(&b.Object)
					isTrackFavourited.On(func(value bool) bool {
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

					isTrackFavourited.Notify(func(oldValue bool) bool {
						if oldValue {
							err := tidal.V1.Favourites.RemoveTrack(context.Background(), secrets.UserID(), trackId)
							if err != nil {
								logger.Error("error while removing track from favourites", "error", err)
								return oldValue
							}
						} else {
							err := tidal.V1.Favourites.AddTrack(context.Background(), secrets.UserID(), trackId)
							if err != nil {
								logger.Error("error while adding track to favourites", "error", err)
								return oldValue
							}
						}

						return !oldValue
					})
				}).
				HAlign(gtk.AlignCenterValue).
				VAlign(gtk.AlignCenterValue).
				WithCSSClass("flat"),
			Button().
				TooltipText(gettext.Get("Add to Queue")).
				IconName("plus-symbolic").
				HAlign(gtk.AlignCenterValue).
				VAlign(gtk.AlignCenterValue).
				ActionName("win.player.queue-track").
				ActionTargetValue(glib.NewVariantString(trackId)).
				WithCSSClass("flat"),
		).
			Margin(10).
			HAlign(gtk.AlignEndValue).
			ToGTK(),
		column,
		0,
		1,
		1,
	)
	return 1
}

func ControlsColumn(track *openapi.Track, grid *gtk.Grid, position int, column int) int {
	if track == nil {
		grid.Attach(
			Box(gtk.OrientationHorizontalValue).ToGTK(),
			column,
			0,
			1,
			1,
		)
		return 1
	}
	return controlsColumn(track.Data.ID, grid, position, column)
}

func LegacyControlsColumn(track *v2.TrackItemData, grid *gtk.Grid, position int, column int) int {
	if track == nil {
		grid.Attach(
			Box(gtk.OrientationHorizontalValue).ToGTK(),
			column,
			0,
			1,
			1,
		)
		return 1
	}
	return controlsColumn(strconv.Itoa(track.ID), grid, position, column)
}
