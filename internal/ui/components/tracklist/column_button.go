package tracklist

import (
	"strconv"

	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"github.com/jwijenbergh/puregotk/v4/glib"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func coverButton(trackId int, grid *gtk.Grid, row int, column int) int {
	button := gtk.NewButton()
	button.SetActionName("win.player.play-track")
	button.SetActionTargetValue(glib.NewVariantInt64(int64(trackId)))
	button.SetFocusable(true)
	button.SetFocusOnClick(true)
	grid.Attach(ManagedWidget(&button.Widget).
		CSS(`button:not(:hover) { background-color: transparent; }`).ToGTK(), 0, row, column, 1)
	return 1
}

func ButtonColumn(track *openapi.Track, grid *gtk.Grid, row int, column int) int {
	trackId := -1
	if parsed, err := strconv.Atoi(track.Data.ID); err == nil {
		trackId = parsed
	}
	return coverButton(trackId, grid, row, column)
}

func LegacyButtonColumn(track *v2.TrackItemData, grid *gtk.Grid, row int, column int) int {
	return coverButton(track.ID, grid, row, column)
}
