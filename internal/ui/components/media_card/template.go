package media_card

import (
	"codeberg.org/dergs/tidalwave/internal/resources"
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gtk"
	"github.com/jwijenbergh/puregotk/v4/pango"
)

func SubTitle(text string) schwifty.Label {
	return Label(text).
		FontSize(14).
		FontWeight(400).
		MaxWidthChars(15).
		Color("#939393").
		HAlign(gtk.AlignStartValue).
		Ellipsis(pango.EllipsizeEndValue)
}

func Card[T any](title string, subTitle schwifty.Widgetable[T], coverUrl string) schwifty.Button {
	return Button().
		Child(
			VStack(
				AspectFrame(
					Image().
						PixelSize(172).
						FromPaintable(resources.MissingAlbum()).
						ConnectConstruct(func(i *gtk.Image) {
							injector.MustInject[*imgutil.ImgUtil]().LoadIntoImage(coverUrl, i)
						}),
				).CornerRadius(10).Overflow(gtk.OverflowHiddenValue),
				Label(title).
					FontSize(16).
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
		WithCSSClass("transparent")
}
