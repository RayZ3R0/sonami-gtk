package syntax

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"codeberg.org/dergs/tonearm/pkg/schwifty/callback"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func Scale(orientation gtk.Orientation) schwifty.Scale {
	return managed("Scale", func() *gtk.Scale {
		scale := gtk.NewScale(orientation, nil)
		scale.ConnectChangeValue(&callback.RangeChangeValueCallback)
		return scale
	})
}
