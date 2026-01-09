package tracklist

import (
	"strconv"

	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func customButtonColumn(trackId string, grid *gtk.Grid, position int, column int, onClick func(trackId string)) int {
	grid.Attach(
		Button().
			ConnectClicked(func(b gtk.Button) {
				onClick(trackId)
			}).
			WithCSSClass("transparent").
			ToGTK(),
		0,
		0,
		column,
		1,
	)
	return 1
}

func CustomButtonColumn(onClick func(trackId string)) ColumnFunc[*openapi.Track] {
	return func(track *openapi.Track, grid *gtk.Grid, position int, column int) int {
		return customButtonColumn(track.Data.ID, grid, position, column, onClick)
	}
}

func LegacyCustomButtonColumn(onClick func(trackId string)) ColumnFunc[*v2.TrackItemData] {
	return func(track *v2.TrackItemData, grid *gtk.Grid, position int, column int) int {
		return customButtonColumn(strconv.Itoa(track.ID), grid, position, column, onClick)
	}
}

func ExpandCustomButtonColumn(additionalWidth int, onClick func(trackId string)) ColumnFunc[*openapi.Track] {
	return func(track *openapi.Track, grid *gtk.Grid, position, column int) int {
		return CustomButtonColumn(onClick)(track, grid, position, column+additionalWidth)
	}
}

func LegacyExpandCustomButtonColumn(additionalWidth int, onClick func(trackId string)) ColumnFunc[*v2.TrackItemData] {
	return func(track *v2.TrackItemData, grid *gtk.Grid, position int, column int) int {
		return LegacyCustomButtonColumn(onClick)(track, grid, position, column+additionalWidth)
	}
}
