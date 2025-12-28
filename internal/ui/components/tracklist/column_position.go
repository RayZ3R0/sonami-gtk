package tracklist

import (
	"strconv"

	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func positionColumn(grid *gtk.Grid, row int, column int) int {
	widget := Label(strconv.Itoa(row + 1)).
		FontWeight(500).
		HAlign(gtk.AlignStartValue).
		HExpand(false).
		Margin(10)
	grid.Attach(widget.ToGTK(), column, row, 1, 1)
	return 1
}

func PositionColumn(track *openapi.Track, grid *gtk.Grid, row int, column int) int {
	return positionColumn(grid, row, column)
}

func LegacyPositionColumn(track *v2.TrackItemData, grid *gtk.Grid, row int, column int) int {
	return positionColumn(grid, row, column)
}
