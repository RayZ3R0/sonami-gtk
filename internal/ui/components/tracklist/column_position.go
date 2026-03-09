package tracklist

import (
	"strconv"

	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
	"codeberg.org/puregotk/puregotk/v4/gtk"
)

func PositionColumn(track sonami.Track, grid *gtk.Grid, position int, column int32) int {
	widget := Label(strconv.Itoa(position + 1)).
		FontWeight(500).
		HAlign(gtk.AlignStartValue).
		HExpand(false).
		Margin(10).HPadding(10)
	grid.Attach(widget.ToGTK(), column, 0, 1, 1)
	return 1
}
