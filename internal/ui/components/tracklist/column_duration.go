package tracklist

import (
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func DurationColumn(track tonearm.Track, grid *gtk.Grid, position int, column int) int {
	grid.Attach(Label(tidalapi.FormatDuration(track.Duration())).Margin(10).ToGTK(), column, 0, 1, 1)
	return 1
}
