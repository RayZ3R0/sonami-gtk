package tracklist

import (
	"codeberg.org/dergs/tonearm/internal/resources"
	"codeberg.org/dergs/tonearm/internal/settings"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
	"codeberg.org/dergs/tonearm/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func CoverColumn(track tonearm.Track, grid *gtk.Grid, position int, column int) int {
	grid.Attach(
		Image().
			FromPaintable(resources.MissingAlbum()).
			PixelSize(54).
			HExpand(false).
			VExpand(false).
			CornerRadius(10).
			Margin(10).
			HAlign(gtk.AlignStartValue).
			Overflow(gtk.OverflowHiddenValue).
			ConnectConstruct(func(i *gtk.Image) {
				coverURL := track.Cover(80)
				if settings.Performance().AllowTracklistImages() && coverURL != "" {
					injector.MustInject[*imgutil.ImgUtil]().LoadIntoImage(coverURL, i)
				}
			}).ToGTK(),
		column,
		0,
		1,
		1,
	)
	return 1
}
