package tracklist_header

import (
	"github.com/RayZ3R0/sonami-gtk/internal/resources"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	"github.com/RayZ3R0/sonami-gtk/pkg/utils/imgutil"
	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/infinytum/injector"
)

func componentCover(coverUrl string) schwifty.AspectFrame {
	return AspectFrame(
		Image().
			PixelSize(146).
			FromPaintable(resources.MissingAlbum()).
			ConnectRealize(func(w gtk.Widget) {
				if coverUrl != "" {
					injector.MustInject[*imgutil.ImgUtil]().LoadIntoImage(coverUrl, gtk.ImageNewFromInternalPtr(w.Ptr))
				}
			}),
	).CornerRadius(10).Overflow(gtk.OverflowHiddenValue).WithCSSClass("cover")
}
