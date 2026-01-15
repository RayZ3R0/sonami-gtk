package tracklist

import (
	"strconv"

	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tonearm/pkg/tidalapi/models/v2"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

type buttonColumnCallback func(trackID string, position, column int)

func customButtonColumn(trackId string, grid *gtk.Grid, position int, column int, onClick buttonColumnCallback) int {
	grid.Attach(
		Button().
			ConnectClicked(func(b gtk.Button) {
				onClick(trackId, position, column)
			}).
			WithCSSClass("transparent").
			ToGTK(),
		0,
		0,
		column,
		1,
	)
	return 0
}

func customWidgetButtonColumn(trackId string, grid *gtk.Grid, position int, column int, button func(trackID string, position, column int) *gtk.Widget) int {
	grid.Attach(
		HStack(
			button(trackId, position, column),
		).
			HAlign(gtk.AlignCenterValue).
			VAlign(gtk.AlignCenterValue).
			ToGTK(),
		column,
		0,
		1,
		1,
	)
	return 0
}

func CustomButtonColumn(onClick buttonColumnCallback) ColumnFunc[*openapi.Track] {
	return func(track *openapi.Track, grid *gtk.Grid, position int, column int) int {
		if track == nil {
			grid.Attach(
				Box(gtk.OrientationHorizontalValue).ToGTK(),
				0,
				0,
				column,
				1,
			)
			return 1
		}
		return customButtonColumn(track.Data.ID, grid, position, column, onClick)
	}
}

func CustomWidgetButtonColumn(button func(string, int, int) *gtk.Widget) ColumnFunc[*openapi.Track] {
	return func(track *openapi.Track, grid *gtk.Grid, position int, column int) int {
		if track == nil {
			grid.Attach(
				Box(gtk.OrientationHorizontalValue).ToGTK(),
				column,
				0,
				1,
				1,
			)
			return 1
		}
		return customWidgetButtonColumn(track.Data.ID, grid, position, column, button)
	}
}

func LegacyCustomButtonColumn(onClick buttonColumnCallback) ColumnFunc[*v2.TrackItemData] {
	return func(track *v2.TrackItemData, grid *gtk.Grid, position int, column int) int {
		if track == nil {
			grid.Attach(
				Box(gtk.OrientationHorizontalValue).ToGTK(),
				0,
				0,
				column,
				1,
			)
			return 1
		}
		return customButtonColumn(strconv.Itoa(track.ID), grid, position, column, onClick)
	}
}

func ExpandCustomButtonColumn(additionalWidth int, onClick buttonColumnCallback) ColumnFunc[*openapi.Track] {
	return func(track *openapi.Track, grid *gtk.Grid, position, column int) int {
		return CustomButtonColumn(onClick)(track, grid, position, column+additionalWidth)
	}
}

func LegacyExpandCustomButtonColumn(additionalWidth int, onClick buttonColumnCallback) ColumnFunc[*v2.TrackItemData] {
	return func(track *v2.TrackItemData, grid *gtk.Grid, position int, column int) int {
		return LegacyCustomButtonColumn(onClick)(track, grid, position, column+additionalWidth)
	}
}
