package tracklist

import (
	"strconv"

	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"github.com/jwijenbergh/puregotk/v4/glib"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func buttonColumn(trackId string, grid *gtk.Grid, position int, column int) int {
	grid.Attach(
		Button().
			ActionName("win.player.play-track").
			ActionTargetValue(glib.NewVariantString(trackId)).
			WithCSSClass("transparent").
			ToGTK(),
		0,
		0,
		column,
		1,
	)
	return 1
}

func ButtonColumn(track *openapi.Track, grid *gtk.Grid, position int, column int) int {
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
	return buttonColumn(track.Data.ID, grid, position, column)
}

func LegacyButtonColumn(track *v2.TrackItemData, grid *gtk.Grid, position int, column int) int {
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
	return buttonColumn(strconv.Itoa(track.ID), grid, position, column)
}

func ExpandButtonColumn(additionalWidth int) ColumnFunc[*openapi.Track] {
	return func(track *openapi.Track, grid *gtk.Grid, position, column int) int {
		return ButtonColumn(track, grid, position, column+additionalWidth)
	}
}

func LegacyExpandButtonColumn(additionalWidth int) ColumnFunc[*v2.TrackItemData] {
	return func(track *v2.TrackItemData, grid *gtk.Grid, position, column int) int {
		return LegacyButtonColumn(track, grid, position, column+additionalWidth)
	}
}
