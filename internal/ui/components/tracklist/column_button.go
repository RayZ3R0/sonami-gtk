package tracklist

import (
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
	"github.com/jwijenbergh/puregotk/v4/glib"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func ButtonColumn(track tonearm.Track, grid *gtk.Grid, position int, column int) int {
	grid.Attach(
		Button().
			ActionName("win.player.play-track").
			ActionTargetValue(glib.NewVariantString(track.ID())).
			WithCSSClass("flat").
			ToGTK(),
		0,
		0,
		column,
		1,
	)
	return 0
}

func ExpandButtonColumn(additionalWidth int) ColumnFunc {
	return func(track tonearm.Track, grid *gtk.Grid, position, column int) int {
		return ButtonColumn(track, grid, position, column+additionalWidth)
	}
}
