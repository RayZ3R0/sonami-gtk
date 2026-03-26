package shortcut_list

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

func NewShortcut(title string, subtitle string, coverUrl string) schwifty.Button {
	return Button().
		Child(
			HStack(
				VStack(
					Label(title).HAlign(gtk.AlignStartValue).WithCSSClass("heading").Ellipsis(pango.EllipsizeEndValue),
					Label(subtitle).HAlign(gtk.AlignStartValue).Visible(subtitle != "").FontWeight(500).WithCSSClass("dimmed").Ellipsis(pango.EllipsizeEndValue),
				).HAlign(gtk.AlignStartValue).VAlign(gtk.AlignCenterValue).HExpand(true),
				AspectFrame(
					Image().
						PixelSize(54).
						FromPaintable(resources.MissingAlbum()).
						ConnectRealize(func(i gtk.Widget) {
							if settings.Performance().AllowShortcutImages() {
								injector.MustInject[*imgutil.ImgUtil]().LoadIntoImageCropped(coverUrl, gtk.ImageNewFromInternalPtr(i.Ptr))
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
					Label(title).HAlign(gtk.AlignCenterValue).WithCSSClass("heading").Ellipsis(pango.EllipsizeEndValue),
					Label(subtitle).HAlign(gtk.AlignCenterValue).Visible(subtitle != "").FontWeight(500).WithCSSClass("dimmed").Ellipsis(pango.EllipsizeEndValue),
				).HAlign(gtk.AlignCenterValue).VAlign(gtk.AlignCenterValue).HExpand(true),
			),
		).
		HExpand(true).
		PaddingStart(15).PaddingEnd(15).VPadding(5)
}
