package tracklist

import (
	"log/slog"
	"strconv"

	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/state"
	favouritebutton "codeberg.org/dergs/tonearm/internal/ui/components/favourite_button"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tonearm/pkg/tidalapi/models/v2"
	"github.com/jwijenbergh/puregotk/v4/glib"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var logger = slog.With("module", "components/tracklist")

func controlsColumn(trackId string, grid *gtk.Grid, position int, column int) int {
	grid.Attach(
		HStack(
			favouritebutton.FavouriteButton(state.TracksCache, trackId).
				HAlign(gtk.AlignCenterValue).
				VAlign(gtk.AlignCenterValue),
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
