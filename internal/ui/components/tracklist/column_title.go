package tracklist

import (
	"codeberg.org/dergs/tidalwave/pkg/gui"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotk4/pkg/pango"
)

func titleColumn(title string, grid *gtk.Grid, row int, column int) int {
	widget := gui.Text(title).FontWeight(500).Ellipsis(pango.EllipsizeEnd).HAlign(gtk.AlignStart).HExpand(true).Margin(10)
	grid.Attach(widget, column, row, 1, 1)
	return 1
}

func TitleColumn(track *openapi.Track, grid *gtk.Grid, row int, column int) int {
	return titleColumn(track.Data.Attributes.Title, grid, row, column)
}

func LegacyTitleColumn(track *v2.TrackItemData, grid *gtk.Grid, row int, column int) int {
	return titleColumn(track.Title, grid, row, column)
}
