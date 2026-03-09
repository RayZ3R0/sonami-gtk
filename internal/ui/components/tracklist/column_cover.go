package tracklist

import (
	"github.com/RayZ3R0/sonami-gtk/internal/resources"
	"github.com/RayZ3R0/sonami-gtk/internal/settings"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
	"github.com/RayZ3R0/sonami-gtk/pkg/utils/imgutil"
	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/infinytum/injector"
)

func CoverColumn(track sonami.Track, grid *gtk.Grid, position int, column int32) int {
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
			ConnectRealize(func(i gtk.Widget) {
				coverURL := track.Cover(80)
				if settings.Performance().AllowTracklistImages() && coverURL != "" {
					injector.MustInject[*imgutil.ImgUtil]().LoadIntoImage(coverURL, gtk.ImageNewFromInternalPtr(i.Ptr))
				}
			}).ToGTK(),
		column,
		0,
		1,
		1,
	)
	return 1
}
