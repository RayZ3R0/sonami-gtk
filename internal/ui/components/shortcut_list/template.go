package shortcut_list

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
)

func NewShortcut(title string, subtitle string, coverUrl string) schwifty.Button {
	// TODO: This might be a memory leak. Ideally we just load this "missing" picture once globally and then share it instead of doing this.
	defaultPixbuf, _ := gdkpixbuf.NewPixbufFromResource("/org/codeberg/dergs/tidalwave/icons/scalable/state/missing-album.svg")
	coverState := state.NewStateful[*gdkpixbuf.Pixbuf](defaultPixbuf)
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
			HStack(
				VStack(
					Label(title).HAlign(gtk.AlignStartValue).FontWeight(600),
					Label(subtitle).HAlign(gtk.AlignStartValue).Visible(subtitle != ""),
				).HAlign(gtk.AlignStartValue).VAlign(gtk.AlignCenterValue).HExpand(true),
				AspectFrame(
					Image().
						PixelSize(54).
						BindPixbuf(coverState),
				).CornerRadius(10).Overflow(gtk.OverflowHiddenValue).HAlign(gtk.AlignEndValue).MarginStart(10),
			),
		).
		HExpand(true).
		MinWidth(300).
		PaddingStart(15).PaddingEnd(5).VPadding(5)
}
