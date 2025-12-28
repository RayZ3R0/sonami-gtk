package syntax

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func Scale(orientation gtk.Orientation) schwifty.Scale {
	return managed("Scale", func() *gtk.Scale {
		return gtk.NewScale(orientation, nil)
	})
}
