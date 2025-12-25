package tracklist

import (
	"context"

	"codeberg.org/dergs/tidalwave/pkg/gui"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotkit/gtkutil/imgutil"
	"github.com/infinytum/injector"
)

func coverColumn(url string, grid *gtk.Grid, row int, column int) int {
	cover := gtk.NewImageFromResource("/org/codeberg/dergs/tidalwave/icons/scalable/state/missing-album.svg")
	cover.SetPixelSize(54)
	cover.SetHExpand(false)
	cover.SetVExpand(false)

	if url != "" {
		imgutil.AsyncGET(
			injector.MustInject[context.Context](),
			url,
			imgutil.ImageSetterFromImage(cover),
		)
	}

	frame := gui.AspectFrame(gui.Wrapper(cover)).
		CornerRadius(10).
		Margin(10).
		HAlign(gtk.AlignStart).
		Overflow(gtk.OverflowHidden)
	grid.Attach(frame, column, row, 1, 1)
	return 1
}

func CoverColumn(track *openapi.Track, grid *gtk.Grid, row int, column int) int {
	url := ""
	for _, album := range track.Included.Albums(track.Data.Relationships.Albums.Data...) {
		for _, artwork := range album.Included.PlainArtworks(album.Data.Relationships.CoverArt.Data...) {
			url = artwork.Attributes.Files.AtLeast(320).Href
			break
		}
	}
	return coverColumn(url, grid, row, column)
}

func LegacyCoverColumn(track *v2.TrackItemData, grid *gtk.Grid, row int, column int) int {
	return coverColumn(tidalapi.ImageURL(track.Album.Cover), grid, row, column)
}
