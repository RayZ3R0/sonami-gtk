package tracklist

import (
	"strings"

	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"github.com/jwijenbergh/puregotk/v4/gtk"
	"github.com/jwijenbergh/puregotk/v4/pango"
)

func artistsColumn(artists []string, grid *gtk.Grid, row int, column int) int {
	grid.Attach(
		Label(strings.Join(artists, ", ")).
			HAlign(gtk.AlignStartValue).
			VAlign(gtk.AlignCenterValue).
			Margin(10).
			Ellipsis(pango.EllipsizeEndValue).
			HExpand(true).
			ToGTK(),
		column,
		row,
		1,
		1,
	)
	return 1
}

func ArtistsColumn(track *openapi.Track, grid *gtk.Grid, row int, column int) int {
	artists := make([]string, 0)
	for _, artist := range track.Included.PlainArtists(track.Data.Relationships.Artists.Data...) {
		artists = append(artists, artist.Attributes.Name)
	}
	return artistsColumn(artists, grid, row, column)
}

func LegacyArtistsColumn(track *v2.TrackItemData, grid *gtk.Grid, row int, column int) int {
	artists := make([]string, 0)
	for _, artist := range track.Artists {
		artists = append(artists, artist.Name)
	}
	return artistsColumn(artists, grid, row, column)
}
