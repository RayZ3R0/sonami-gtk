package tracklist

import (
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func controlsColumn(grid *gtk.Grid, row int, column int) int {
	grid.Attach(
		HStack(
			Button().
				IconName("heart-outline-thick-symbolic").
				HAlign(gtk.AlignCenterValue).
				VAlign(gtk.AlignCenterValue).
				CSS(`button:not(:hover) { background-color: transparent; }`),
			Button().
				IconName("plus-symbolic").
				HAlign(gtk.AlignCenterValue).
				VAlign(gtk.AlignCenterValue).
				CSS(`button:not(:hover) { background-color: transparent; }`),
		).
			Margin(10).
			HAlign(gtk.AlignEndValue).
			ToGTK(),
		column,
		row,
		1,
		1,
	)
	return 1
}

func ControlsColumn(track *openapi.Track, grid *gtk.Grid, row int, column int) int {
	return controlsColumn(grid, row, column)
}

func LegacyControlsColumn(track *v2.TrackItemData, grid *gtk.Grid, row int, column int) int {
	return controlsColumn(grid, row, column)
}
