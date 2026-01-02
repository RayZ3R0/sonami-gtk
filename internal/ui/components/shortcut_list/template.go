package shortcut_list

import (
	"codeberg.org/dergs/tidalwave/internal/resources"
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func NewShortcut(title string, subtitle string, coverUrl string) schwifty.Button {
	return Button().
		Child(
			HStack(
				VStack(
					Label(title).HAlign(gtk.AlignStartValue).FontWeight(600).FontSize(16),
					Label(subtitle).HAlign(gtk.AlignStartValue).Visible(subtitle != "").FontSize(14).FontWeight(500).Color("#939393"),
				).HAlign(gtk.AlignStartValue).VAlign(gtk.AlignCenterValue).HExpand(true),
				AspectFrame(
					Image().
						PixelSize(54).
						FromPaintable(resources.MissingAlbum()).
						ConnectConstruct(func(i *gtk.Image) {
							injector.MustInject[*imgutil.ImgUtil]().LoadIntoImage(coverUrl, i)
						}),
				).CornerRadius(10).Overflow(gtk.OverflowHiddenValue).HAlign(gtk.AlignEndValue).MarginStart(10),
			),
		).
		HExpand(true).
		MinWidth(300).
		PaddingStart(15).PaddingEnd(5).VPadding(5)
}
