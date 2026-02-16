package tracklist

import (
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func groupedColumn(columns []any, width int, grid *gtk.Grid, row int, column int) int {
	grid.Attach(
		HStack(columns...).
			ToGTK(),
		column,
		row,
		width,
		1,
	)
	return width
}

func GroupedColumn(width int, align gtk.Align, columns ...ColumnFunc) ColumnFunc {
	return func(track tonearm.Track, grid *gtk.Grid, position, column int) int {
		subGrid := gtk.NewGrid()
		subGrid.SetValign(gtk.AlignCenterValue)
		subGrid.SetHalign(align)
		subWidth := 0
		for i, c := range columns {
			subWidth += c(track, subGrid, position, column+i)
		}
		defer subGrid.Unref()
		grid.Attach(
			&subGrid.Widget,
			column,
			0,
			width,
			1,
		)
		return width
	}
}
