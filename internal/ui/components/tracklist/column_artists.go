package tracklist

import (
	"strings"

	"codeberg.org/dergs/tidalwave/pkg/gui"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotk4/pkg/pango"
)

func artistsColumn(artists []string, grid *gtk.Grid, row int, column int) int {
	widget := gui.Text(strings.Join(artists, ", ")).
		HAlign(gtk.AlignStart).
		Margin(10).
		Ellipsis(pango.EllipsizeEnd).
		HExpand(true)
	grid.Attach(widget, column, row, 1, 1)
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
