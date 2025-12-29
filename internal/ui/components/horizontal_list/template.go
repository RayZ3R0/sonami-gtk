package horizontal_list

import (
	"log/slog"

	"codeberg.org/dergs/tidalwave/internal/g"
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gdkpixbuf"
	"github.com/jwijenbergh/puregotk/v4/glib"
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
	coverState := state.NewStateful[*gdkpixbuf.Pixbuf](nil)
	if coverUrl != "" {
		go func() {
			pixbuf, err := injector.MustInject[*imgutil.ImgUtil]().LoadPixbuf(coverUrl)
			if err != nil {
				slog.Error("failed to load album cover", "error", err)
				return
			}
			glib.IdleAddOnce(
				g.Ptr[glib.SourceOnceFunc](func(u uintptr) {
					coverState.SetValue(pixbuf)
					pixbuf.Unref()
				}),
				0,
			)
		}()
	}

	return Button().
		Child(
			VStack(
				AspectFrame(
					Image().
						PixelSize(172).
						FromResource("/org/codeberg/dergs/tidalwave/icons/scalable/state/missing-album.svg").
						BindPixbuf(coverState),
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
		CSS("button:not(:hover) { background-color: transparent; }")
}
