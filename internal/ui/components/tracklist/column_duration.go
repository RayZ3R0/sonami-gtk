package tracklist

import (
	"time"

	"codeberg.org/dergs/tidalwave/pkg/gui"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func durationColumn(duration time.Duration, grid *gtk.Grid, row int, column int) int {
	grid.Attach(gui.Text(tidalapi.FormatDuration(int(duration.Seconds()))).Margin(10), column, row, 1, 1)
	return 1
}

func DurationColumn(track *openapi.Track, grid *gtk.Grid, row int, column int) int {
	return durationColumn(track.Data.Attributes.Duration.Duration, grid, row, column)
}

func LegacyDurationColumn(track *v2.TrackItemData, grid *gtk.Grid, row int, column int) int {
	return durationColumn(time.Duration(track.Duration)*time.Second, grid, row, column)
}
