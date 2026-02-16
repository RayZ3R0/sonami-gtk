package tracklist

import (
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
	"github.com/jwijenbergh/puregotk/v4/gtk"
	"github.com/jwijenbergh/puregotk/v4/pango"
)

func TitleAlbumColumn(track tonearm.Track, grid *gtk.Grid, position int, column int) int {
	frame := VStack(
		Label(track.Title()).FontWeight(500).Ellipsis(pango.EllipsizeEndValue).HAlign(gtk.AlignStartValue),
		Label(track.Album().Title()).WithCSSClass("dimmed").Ellipsis(pango.EllipsizeEndValue).HAlign(gtk.AlignStartValue),
	).Spacing(3).VAlign(gtk.AlignCenterValue).HAlign(gtk.AlignStartValue).HExpand(true).Margin(10)
	grid.Attach(HStack(frame, Spacer()).ToGTK(), column, 0, 1, 1)
	return 1
}
