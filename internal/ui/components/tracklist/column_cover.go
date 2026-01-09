package tracklist

import (
	"codeberg.org/dergs/tidalwave/internal/resources"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"codeberg.org/dergs/tidalwave/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func coverColumn(url string, grid *gtk.Grid, position int, column int) int {
	grid.Attach(
		AspectFrame(
			Image().
				FromPaintable(resources.MissingAlbum()).
				PixelSize(54).
				HExpand(false).
				VExpand(false).
				ConnectConstruct(func(i *gtk.Image) {
					injector.MustInject[*imgutil.ImgUtil]().LoadIntoImage(url, i)
				}),
		).
			CornerRadius(10).
			Margin(10).
			HAlign(gtk.AlignStartValue).
			Overflow(gtk.OverflowHiddenValue).ToGTK(),
		column,
		0,
		1,
		1,
	)
	return 1
}

func CoverColumn(track *openapi.Track, grid *gtk.Grid, position int, column int) int {
	url := ""
	for _, album := range track.Included.Albums(track.Data.Relationships.Albums.Data...) {
		for _, artwork := range album.Included.PlainArtworks(album.Data.Relationships.CoverArt.Data...) {
			url = artwork.Attributes.Files.AtLeast(80).Href
			break
		}
	}
	return coverColumn(url, grid, position, column)
}

func LegacyCoverColumn(track *v2.TrackItemData, grid *gtk.Grid, position int, column int) int {
	return coverColumn(tidalapi.ImageURL(track.Album.Cover), grid, position, column)
}
