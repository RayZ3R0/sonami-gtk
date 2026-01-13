package syntax

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func CenterBox() schwifty.CenterBox {
	return managed("CenterBox", func() *gtk.CenterBox {
		return gtk.NewCenterBox()
	})
}
