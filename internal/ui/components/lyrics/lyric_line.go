package lyrics

import (
	"time"

	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

type lyricTiming struct {
	timed              bool
	timeStart, timeEnd time.Duration
}

func lyricLine(text string, timing lyricTiming) schwifty.Button {
	return Button().
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
