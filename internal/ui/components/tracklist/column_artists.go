package tracklist

import (
	"strings"

	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tonearm/pkg/tidalapi/models/v2"
	"github.com/jwijenbergh/puregotk/v4/gtk"
	"github.com/jwijenbergh/puregotk/v4/pango"
)

func artistsColumn(artists []string, grid *gtk.Grid, position int, column int) int {
	grid.Attach(
		Label(strings.Join(artists, ", ")).
			HAlign(gtk.AlignStartValue).
			VAlign(gtk.AlignCenterValue).
			Margin(10).
			Ellipsis(pango.EllipsizeEndValue).
			HExpand(true).
			ToGTK(),
		column,
		0,
		1,
		1,
	)
	return 1
}

func ArtistsColumn(track *openapi.Track, grid *gtk.Grid, position int, column int) int {
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
	artists := make([]string, 0)
	for _, artist := range track.Included.PlainArtists(track.Data.Relationships.Artists.Data...) {
		artists = append(artists, artist.Attributes.Name)
	}
	return artistsColumn(artists, grid, position, column)
}

func LegacyArtistsColumn(track *v2.TrackItemData, grid *gtk.Grid, position int, column int) int {
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

	artists := make([]string, 0)
	for _, artist := range track.Artists {
		artists = append(artists, artist.Name)
	}
	return artistsColumn(artists, grid, position, column)
}
