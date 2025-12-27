package tracklist

import (
	"strconv"

	"codeberg.org/dergs/tidalwave/pkg/gui"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func positionColumn(grid *gtk.Grid, row int, column int) int {
	widget := gui.Text(strconv.Itoa(row + 1)).FontWeight(500).HAlign(gtk.AlignStart).HExpand(false).Margin(10)
	grid.Attach(widget, column, row, 1, 1)
	return 1
}

func PositionColumn(track *openapi.Track, grid *gtk.Grid, row int, column int) int {
	return positionColumn(grid, row, column)
}

func LegacyPositionColumn(track *v2.TrackItemData, grid *gtk.Grid, row int, column int) int {
	return positionColumn(grid, row, column)
}
