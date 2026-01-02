package syntax

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func ScrolledWindow() schwifty.ScrolledWindow {
	return managed("ScrolledWindow", func() *gtk.ScrolledWindow {
		return gtk.NewScrolledWindow()
	})
}
