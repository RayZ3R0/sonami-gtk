package components

import (
	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
)

func SquirclePicture(picture schwifty.Picture) schwifty.Picture {
	return picture.
		Background("alpha(var(--view-fg-color), 0.1)").
		Overflow(gtk.OverflowHiddenValue).CornerRadius(10)
}
