package tracklist

import (
	"strconv"

	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"github.com/jwijenbergh/puregotk/v4/glib"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func buttonColumn(trackId int, grid *gtk.Grid, row int, column int) int {
	grid.Attach(
		Button().
			ActionName("win.player.play-track").
			ActionTargetValue(glib.NewVariantInt64(int64(trackId))).
			CSS(`button:not(:hover) { background-color: transparent; }`).
			ToGTK(),
		0,
		row,
		column,
		1,
	)
	return 1
}

func ButtonColumn(track *openapi.Track, grid *gtk.Grid, row int, column int) int {
	trackId := -1
	if parsed, err := strconv.Atoi(track.Data.ID); err == nil {
		trackId = parsed
	}
	return buttonColumn(trackId, grid, row, column)
}

func LegacyButtonColumn(track *v2.TrackItemData, grid *gtk.Grid, row int, column int) int {
	return buttonColumn(track.ID, grid, row, column)
}
