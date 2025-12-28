package tracklist

import (
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"codeberg.org/dergs/tidalwave/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func coverColumn(url string, grid *gtk.Grid, row int, column int) int {
	cover := gtk.NewImageFromResource("/org/codeberg/dergs/tidalwave/icons/scalable/state/missing-album.svg")
	cover.SetPixelSize(54)
	cover.SetHexpand(false)
	cover.SetVexpand(false)
	defer cover.Unref()

	aspectFrame := gtk.NewAspectFrame(0.5, 0.5, 1.0, false)
	aspectFrame.SetChild(&cover.Widget)

	if url != "" {
		injector.MustInject[*imgutil.ImgUtil]().LoadIntoImage(url, cover)
	}

	frame := ManagedWidget(&aspectFrame.Widget).
		CornerRadius(10).
		Margin(10).
		HAlign(gtk.AlignStartValue).
		Overflow(gtk.OverflowHiddenValue)
	grid.Attach(frame.ToGTK(), column, row, 1, 1)
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
