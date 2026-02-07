package shortcut_list

import (
	"codeberg.org/dergs/tonearm/internal/resources"
	"codeberg.org/dergs/tonearm/internal/settings"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func NewShortcut(title string, subtitle string, coverUrl string) schwifty.Button {
	return Button().
		Child(
			HStack(
				VStack(
					Label(title).HAlign(gtk.AlignStartValue).WithCSSClass("heading"),
					Label(subtitle).HAlign(gtk.AlignStartValue).Visible(subtitle != "").FontWeight(500).WithCSSClass("dimmed"),
				).HAlign(gtk.AlignStartValue).VAlign(gtk.AlignCenterValue).HExpand(true),
				AspectFrame(
					Image().
						PixelSize(54).
						FromPaintable(resources.MissingAlbum()).
						ConnectConstruct(func(i *gtk.Image) {
							if settings.Performance().AllowShortcutImages() {
								injector.MustInject[*imgutil.ImgUtil]().LoadIntoImage(coverUrl, i)
							}
						}),
				).CornerRadius(10).Overflow(gtk.OverflowHiddenValue).HAlign(gtk.AlignEndValue).MarginStart(10),
			),
		).
		HExpand(true).
		MinWidth(300).
		PaddingStart(15).PaddingEnd(5).VPadding(5)
}

func NewTextShortcut(title string, subtitle string) schwifty.Button {
	return Button().
		Child(
			HStack(
				VStack(
					Label(title).HAlign(gtk.AlignCenterValue).WithCSSClass("heading"),
					Label(subtitle).HAlign(gtk.AlignCenterValue).Visible(subtitle != "").FontWeight(500).WithCSSClass("dimmed"),
				).HAlign(gtk.AlignCenterValue).VAlign(gtk.AlignCenterValue).HExpand(true),
			),
		).
		HExpand(true).
		PaddingStart(15).PaddingEnd(15).VPadding(5)
}
