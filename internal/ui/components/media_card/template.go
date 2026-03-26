package media_card

import (
	"codeberg.org/puregotk/puregotk/v4/gtk"
	"codeberg.org/puregotk/puregotk/v4/pango"
	"github.com/RayZ3R0/sonami-gtk/internal/resources"
	"github.com/RayZ3R0/sonami-gtk/internal/settings"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	"github.com/RayZ3R0/sonami-gtk/pkg/utils/imgutil"
	"github.com/infinytum/injector"
)

func SubTitle(text string) schwifty.Label {
	return Label(text).
		FontWeight(400).
		MaxWidthChars(15).
		WithCSSClass("dimmed").
		HAlign(gtk.AlignStartValue).
		Ellipsis(pango.EllipsizeEndValue)
}

func Card[T any](title string, subTitle schwifty.Widgetable[T], coverUrl string) schwifty.Button {
	return Button().
		Child(
			VStack(
				Image().
					PixelSize(172).
					FromPaintable(resources.MissingAlbum()).
					ConnectRealize(func(i gtk.Widget) {
						if settings.Performance().AllowMediaCardImages() {
							injector.MustInject[*imgutil.ImgUtil]().LoadIntoImageCropped(coverUrl, gtk.ImageNewFromInternalPtr(i.Ptr))
						}
					}).CornerRadius(10).Overflow(gtk.OverflowHiddenValue),
				Label(title).
					WithCSSClass("heading").
					MarginTop(10).
					MaxWidthChars(15).
					HAlign(gtk.AlignStartValue).
					Ellipsis(pango.EllipsizeEndValue),
				subTitle.MarginTop(2),
			),
		).
		Padding(10).
		HExpand(false).
		VExpand(false).
		WithCSSClass("flat")
}
