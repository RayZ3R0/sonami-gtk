package tracklist

import (
	"strconv"

	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"github.com/jwijenbergh/puregotk/v4/glib"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func controlsColumn(trackId string, grid *gtk.Grid, position int, column int) int {
	grid.Attach(
		HStack(
			Button().
				IconName("heart-outline-thick-symbolic").
				HAlign(gtk.AlignCenterValue).
				VAlign(gtk.AlignCenterValue).
				WithCSSClass("transparent"),
			Button().
				IconName("plus-symbolic").
				HAlign(gtk.AlignCenterValue).
				VAlign(gtk.AlignCenterValue).
				ActionName("win.player.queue-track").
				ActionTargetValue(glib.NewVariantString(trackId)).
				WithCSSClass("transparent"),
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
	return controlsColumn(track.Data.ID, grid, position, column)
}

func LegacyControlsColumn(track *v2.TrackItemData, grid *gtk.Grid, position int, column int) int {
	return controlsColumn(strconv.Itoa(track.ID), grid, position, column)
}
