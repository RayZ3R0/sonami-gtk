package tracklist

import (
	"codeberg.org/dergs/tidalwave/pkg/gui"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func controlsColumn(grid *gtk.Grid, row int, column int) int {
	addToCollection := gtk.NewButtonFromIconName("heart-outline-thick-symbolic")
	addToQueue := gtk.NewButtonFromIconName("plus-symbolic")
	widget := gui.HStack(
		gui.Wrapper(addToCollection).
			HAlign(gtk.AlignCenter).
			VAlign(gtk.AlignCenter).
			CSS(`button:not(:hover) { background-color: transparent; }`),
		gui.Wrapper(addToQueue).
			HAlign(gtk.AlignCenter).
			VAlign(gtk.AlignCenter).
			CSS(`button:not(:hover) { background-color: transparent; }`),
	)
	grid.Attach(widget.Margin(10).HAlign(gtk.AlignEnd), column, row, 1, 1)
	return 1
}

func ControlsColumn(track *openapi.Track, grid *gtk.Grid, row int, column int) int {
	return controlsColumn(grid, row, column)
}

func LegacyControlsColumn(track *v2.TrackItemData, grid *gtk.Grid, row int, column int) int {
	return controlsColumn(grid, row, column)
}
