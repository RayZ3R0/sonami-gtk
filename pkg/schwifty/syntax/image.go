package syntax

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func Image() schwifty.Image {
	return managed("Image", func() *gtk.Image {
		return gtk.NewImage()
	})
}
