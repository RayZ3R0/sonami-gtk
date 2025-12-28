package syntax

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func Label(text string) schwifty.Label {
	return managed("Label", func() *gtk.Label {
		return gtk.NewLabel(text)
	})
}
