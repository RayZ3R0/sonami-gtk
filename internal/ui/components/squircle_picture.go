package components

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func SquirclePicture(picture schwifty.Picture) schwifty.Picture {
	return picture.
		Background("alpha(var(--view-fg-color), 0.1)").
		Overflow(gtk.OverflowHiddenValue).CornerRadius(10)
}
