package tracklist

import (
	"strings"

	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
	"github.com/jwijenbergh/puregotk/v4/gtk"
	"github.com/jwijenbergh/puregotk/v4/pango"
)

func ArtistsColumn(track tonearm.Track, grid *gtk.Grid, position int, column int) int {
	grid.Attach(
		Label(strings.Join(track.Artists().Names(), ", ")).
			HAlign(gtk.AlignStartValue).
			VAlign(gtk.AlignCenterValue).
			Margin(10).
			Ellipsis(pango.EllipsizeEndValue).
			HExpand(true).
			ToGTK(),
		column,
		0,
		1,
		1,
	)
	return 1
}
