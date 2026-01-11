package tracklist

import (
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"github.com/jwijenbergh/puregotk/v4/gtk"
	"github.com/jwijenbergh/puregotk/v4/pango"
)

func titleAlbumColumn(title string, album string, grid *gtk.Grid, position int, column int) int {
	frame := VStack(
		Label(title).FontWeight(500).Ellipsis(pango.EllipsizeEndValue).HAlign(gtk.AlignStartValue),
		Label(album).CSS(`label { color: #939393; }`).Ellipsis(pango.EllipsizeEndValue).HAlign(gtk.AlignStartValue),
	).Spacing(3).VAlign(gtk.AlignCenterValue).HAlign(gtk.AlignStartValue).HExpand(true).Margin(10)
	grid.Attach(HStack(frame, Spacer()).ToGTK(), column, 0, 1, 1)
	return 1
}

func TitleAlbumColumn(track *openapi.Track, grid *gtk.Grid, position int, column int) int {
	if track == nil {
		frame := VStack(
			Label(""),
			Label(""),
		).Spacing(3).VAlign(gtk.AlignCenterValue).HAlign(gtk.AlignStartValue).HExpand(true).Margin(10)
		grid.Attach(HStack(frame, Spacer()).ToGTK(), column, 0, 1, 1)
		return 1
	}
	albumName := ""
	for _, album := range track.Included.PlainAlbums(track.Data.Relationships.Albums.Data...) {
		albumName = album.Attributes.Title
		break
	}
	return titleAlbumColumn(track.Data.Attributes.Title, albumName, grid, position, column)
}

func LegacyTitleAlbumColumn(track *v2.TrackItemData, grid *gtk.Grid, position int, column int) int {
	if track == nil {
		frame := VStack(
			Label(""),
			Label(""),
		).Spacing(3).VAlign(gtk.AlignCenterValue).HAlign(gtk.AlignStartValue).HExpand(true).Margin(10)
		grid.Attach(HStack(frame, Spacer()).ToGTK(), column, 0, 1, 1)
		return 1
	}
	return titleAlbumColumn(track.Title, track.Album.Title, grid, position, column)
}
