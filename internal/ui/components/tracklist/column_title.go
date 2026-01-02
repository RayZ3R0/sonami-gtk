package tracklist

import (
	"fmt"

	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"github.com/jwijenbergh/puregotk/v4/gtk"
	"github.com/jwijenbergh/puregotk/v4/pango"
)

func titleColumn(title string, grid *gtk.Grid, row int, column int) int {
	widget := Label(title).FontWeight(500).Ellipsis(pango.EllipsizeEndValue).HAlign(gtk.AlignStartValue).HExpand(true).Margin(10)
	grid.Attach(widget.ToGTK(), column, row, 1, 1)
	return 1
}

func TitleColumn(track *openapi.Track, grid *gtk.Grid, row int, column int) int {
	title := track.Data.Attributes.Title
	if version := track.Data.Attributes.Version; version != nil && version != "" {
		title = fmt.Sprintf("%s (%s)", title, version)
	}
	return titleColumn(title, grid, row, column)
}

func LegacyTitleColumn(track *v2.TrackItemData, grid *gtk.Grid, row int, column int) int {
	return titleColumn(track.Title, grid, row, column)
}
