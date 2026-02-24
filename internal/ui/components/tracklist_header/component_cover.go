package tracklist_header

import (
	"codeberg.org/dergs/tonearm/internal/resources"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gtk"
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
