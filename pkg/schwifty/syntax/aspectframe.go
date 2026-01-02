package syntax

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func AspectFrame(child any) schwifty.AspectFrame {
	return managed("Scale", func() *gtk.AspectFrame {
		aspectFrame := gtk.NewAspectFrame(0.5, 0.5, 1.0, false)
		aspectFrame.SetChild(schwifty.ResolveWidget(child))
		return aspectFrame
	})
}
