package tracklist

import (
	"time"

	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tonearm/pkg/tidalapi/models/v2"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func durationColumn(duration time.Duration, grid *gtk.Grid, position int, column int) int {
	grid.Attach(Label(tidalapi.FormatDuration(duration)).Margin(10).ToGTK(), column, 0, 1, 1)
	return 1
}

func DurationColumn(track *openapi.Track, grid *gtk.Grid, position int, column int) int {
	return durationColumn(track.Data.Attributes.Duration.Duration, grid, position, column)
}

func LegacyDurationColumn(track *v2.TrackItemData, grid *gtk.Grid, position int, column int) int {
	return durationColumn(time.Duration(track.Duration)*time.Second, grid, position, column)
}
