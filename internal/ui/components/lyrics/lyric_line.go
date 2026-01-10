package lyrics

import (
	"time"

	"codeberg.org/dergs/tidalwave/internal/player"
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

type lyricTiming struct {
	timeStart, timeEnd time.Duration
}

func lyricLine(text string, timing *lyricTiming) schwifty.Button {
	var classListener string

	b := Button().
		Child(
			lyricLineText(text),
		).
		HExpand(true).
		PaddingTop(24).
		PaddingBottom(24).
		PaddingStart(16).
		PaddingEnd(16).
		CornerRadius(12).
		WithCSSClass("lyric")

	if timing != nil {
		b = b.
			WithCSSClass("timed").
			ConnectConstruct(func(b *gtk.Button) {
				ptr := b.GoPointer()
				classListener = activeLyricIndex.AddCallback(func(newValue uintptr) {
					widget := gtk.ButtonNewFromInternalPtr(ptr)
					schwifty.OnMainThreadOncePure(func() {
						if newValue == ptr {
							widget.AddCssClass("active")
						} else {
							widget.RemoveCssClass("active")
						}
					})
				})
			}).
			ConnectDestroy(func(w gtk.Widget) {
				activeLyricIndex.RemoveCallback(classListener)
			}).
			ConnectClicked(func(gtk.Button) {
				userManuallyScrolled.SetValue(false)
				player.SeekToPosition(timing.timeStart)
			})
	}

	return b
}

func lyricLineText(lyricText string) schwifty.Label {
	return Label(lyricText).
		HAlign(gtk.AlignCenterValue).
		VAlign(gtk.AlignCenterValue).
		FontSize(20).
		FontWeight(600).
		Wrap(true).
		Justify(gtk.JustifyCenterValue)
}
