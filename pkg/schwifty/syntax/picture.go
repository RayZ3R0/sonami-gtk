package syntax

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func Picture() schwifty.Picture {
	return managed("Picture", func() *gtk.Picture {
		return gtk.NewPicture()
	})
}
