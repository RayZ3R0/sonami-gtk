package tracklist

import (
	"codeberg.org/dergs/tidalwave/pkg/gui"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotk4/pkg/pango"
)

func titleAlbumColumn(title string, album string, grid *gtk.Grid, row int, column int) int {
	frame := gui.VStack(
		gui.Text(title).FontWeight(500).Ellipsis(pango.EllipsizeEnd).HAlign(gtk.AlignStart),
		gui.Text(album).CSS(`label { color: #939393; }`).Ellipsis(pango.EllipsizeEnd).HAlign(gtk.AlignStart),
	).Spacing(3).VAlign(gtk.AlignCenter).HAlign(gtk.AlignStart).HExpand(true).Margin(10)
	grid.Attach(gui.HStack(frame, gui.Spacer()), column, row, 1, 1)
	return 1
}

func TitleAlbumColumn(track *openapi.Track, grid *gtk.Grid, row int, column int) int {
	albumName := ""
	for _, album := range track.Included.PlainAlbums(track.Data.Relationships.Albums.Data...) {
		albumName = album.Attributes.Title
		break
	}
	return titleAlbumColumn(track.Data.Attributes.Title, albumName, grid, row, column)
}

func LegacyTitleAlbumColumn(track *v2.TrackItemData, grid *gtk.Grid, row int, column int) int {
	return titleAlbumColumn(track.Title, track.Album.Title, grid, row, column)
}
