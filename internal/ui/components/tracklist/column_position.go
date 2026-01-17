package tracklist

import (
	"strconv"

	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tonearm/pkg/tidalapi/models/v2"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func positionColumn(grid *gtk.Grid, position int, column int) int {
	widget := Label(strconv.Itoa(position + 1)).
		FontWeight(500).
		HAlign(gtk.AlignStartValue).
		HExpand(false).
		Margin(10)
	grid.Attach(widget.ToGTK(), column, 0, 1, 1)
	return 1
}

func PositionColumn(track *openapi.Track, grid *gtk.Grid, position int, column int) int {
	if track == nil {
		widget := Label("").
			FontWeight(500).
			HAlign(gtk.AlignStartValue).
			HExpand(false).
			Margin(10)
		grid.Attach(widget.ToGTK(), column, 0, 1, 1)
		return 1
	}
	return positionColumn(grid, position, column)
}

func LegacyPositionColumn(track *v2.TrackItemData, grid *gtk.Grid, position int, column int) int {
	if track == nil {
		widget := Label("").
			FontWeight(500).
			HAlign(gtk.AlignStartValue).
			HExpand(false).
			Margin(10)
		grid.Attach(widget.ToGTK(), column, 0, 1, 1)
		return 1
	}
	return positionColumn(grid, position, column)
}
