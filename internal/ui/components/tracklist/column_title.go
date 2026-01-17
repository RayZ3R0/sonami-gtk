package tracklist

import (
	"fmt"

	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tonearm/pkg/tidalapi/models/v2"
	"github.com/jwijenbergh/puregotk/v4/gtk"
	"github.com/jwijenbergh/puregotk/v4/pango"
)

func titleColumn(title string, grid *gtk.Grid, position int, column int) int {
	widget := Label(title).FontWeight(500).Ellipsis(pango.EllipsizeEndValue).HAlign(gtk.AlignStartValue).HExpand(true).Margin(10)
	grid.Attach(widget.ToGTK(), column, 0, 1, 1)
	return 1
}

func TitleColumn(track *openapi.Track, grid *gtk.Grid, position int, column int) int {
	if track == nil {
		title := ""
		widget := Label(title).FontWeight(500).Ellipsis(pango.EllipsizeEndValue).HAlign(gtk.AlignStartValue).HExpand(true).Margin(10)
		grid.Attach(widget.ToGTK(), column, 0, 1, 1)
		return 1
	}
	title := track.Data.Attributes.Title
	if version := track.Data.Attributes.Version; version != nil && version != "" {
		title = fmt.Sprintf("%s (%s)", title, version)
	}
	return titleColumn(title, grid, position, column)
}

func LegacyTitleColumn(track *v2.TrackItemData, grid *gtk.Grid, position int, column int) int {
	if track == nil {
		title := ""
		widget := Label(title).FontWeight(500).Ellipsis(pango.EllipsizeEndValue).HAlign(gtk.AlignStartValue).HExpand(true).Margin(10)
		grid.Attach(widget.ToGTK(), column, 0, 1, 1)
		return 1
	}
	return titleColumn(track.Title, grid, position, column)
}
