package tracklist

import (
	"strconv"

	"codeberg.org/dergs/tidalwave/internal/player"
	"codeberg.org/dergs/tidalwave/pkg/gui"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotkit/gtkutil/cssutil"
)

var columnBoxCSS = cssutil.Applier("track-list-column-box", `
	.track-list-column-box {
		border-radius: 10px;
		transition-duration: 0.2s;
		transition-property: background-color;
		transition-timing-function: cubic-bezier(0.25, 0.46, 0.45, 0.94);
	}

	.track-list-column-box:hover {
		background-color: alpha(var(--view-fg-color), 0.15);
	}

	.track-list-column-box:focus:active {
		background-color: alpha(var(--view-fg-color), 0.30);
	}
`)

func coverBox(trackId int, grid *gtk.Grid, row int, column int) int {
	box := gui.HStack().HExpand(true).Focusable(true).FocusOnClick(true)
	columnBoxCSS(box)

	ctrl := gtk.NewGestureClick()
	ctrl.ConnectPressed(func(nPress int, x, y float64) {
		box := gui.Wrapper(ctrl.Widget()).GTKWidget()
		if child := box.FocusChild(); child != nil {
			if gui.Wrapper(child).GTKWidget().StateFlags()&gtk.StateFlagActive != 0 {
				return
			}
		}
		box.GrabFocus()
		box.SetStateFlags(gtk.StateFlagActive, false)
	})
	ctrl.ConnectReleased(func(nPress int, x, y float64) {
		player.Play(trackId)
	})
	box.GTKWidget().AddController(ctrl)

	grid.Attach(box, 0, row, column, 1)
	return 1
}

func BoxColumn(track *openapi.Track, grid *gtk.Grid, row int, column int) int {
	trackId := -1
	if parsed, err := strconv.Atoi(track.Data.ID); err == nil {
		trackId = parsed
	}
	return coverBox(trackId, grid, row, column)
}

func LegacyBoxColumn(track *v2.TrackItemData, grid *gtk.Grid, row int, column int) int {
	return coverBox(track.ID, grid, row, column)
}
