package syntax

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func Button() schwifty.Button {
	return managed("Button", func() *gtk.Button {
		return gtk.NewButton()
	})
}
