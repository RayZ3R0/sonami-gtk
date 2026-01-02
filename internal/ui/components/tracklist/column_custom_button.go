package tracklist

import (
	"strconv"

	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func customButtonColumn(trackId string, grid *gtk.Grid, row int, column int, onClick func(trackId string)) int {
	grid.Attach(
		Button().
			ConnectClicked(func(b gtk.Button) {
				onClick(trackId)
			}).
			WithCSSClass("transparent").
			ToGTK(),
		0,
		row,
		column,
		1,
	)
	return 1
}

func CustomButtonColumn(onClick func(trackId string)) ColumnFunc {
	return func(track *openapi.Track, grid *gtk.Grid, row int, column int) int {
		return customButtonColumn(track.Data.ID, grid, row, column, onClick)
	}
}

func LegacyCustomButtonColumn(onClick func(trackId string)) LegacyColumnFunc {
	return func(track *v2.TrackItemData, grid *gtk.Grid, row int, column int) int {
		return customButtonColumn(strconv.Itoa(track.ID), grid, row, column, onClick)
	}
}
