package tracklist

import (
	"time"

	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func durationColumn(duration time.Duration, grid *gtk.Grid, row int, column int) int {
	grid.Attach(Label(tidalapi.FormatDuration(int(duration.Seconds()))).Margin(10).ToGTK(), column, row, 1, 1)
	return 1
}

func DurationColumn(track *openapi.Track, grid *gtk.Grid, row int, column int) int {
	return durationColumn(track.Data.Attributes.Duration.Duration, grid, row, column)
}

func LegacyDurationColumn(track *v2.TrackItemData, grid *gtk.Grid, row int, column int) int {
	return durationColumn(time.Duration(track.Duration)*time.Second, grid, row, column)
}
