package tracklist

import (
	"codeberg.org/puregotk/puregotk/v4/gtk"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi"
)

func DurationColumn(track sonami.Track, grid *gtk.Grid, position int, column int32) int {
	grid.Attach(Label(tidalapi.FormatDuration(track.Duration())).Margin(10).ToGTK(), column, 0, 1, 1)
	return 1
}
