package tracklist

import (
	"codeberg.org/puregotk/puregotk/v4/gtk"
	"codeberg.org/puregotk/puregotk/v4/pango"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
)

func TitleColumn(track sonami.Track, grid *gtk.Grid, position int, column int32) int {
	widget := Label(sonami.FormatTitle(track)).FontWeight(500).Ellipsis(pango.EllipsizeEndValue).HAlign(gtk.AlignStartValue).HExpand(true).Margin(10)
	grid.Attach(widget.ToGTK(), column, 0, 1, 1)
	return 1
}
