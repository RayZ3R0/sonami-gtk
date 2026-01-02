package syntax

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"github.com/jwijenbergh/puregotk/v4/adw"
)

func WrapBox(children ...any) schwifty.WrapBox {
	return managed("WrapBox", func() *adw.WrapBox {
		box := adw.NewWrapBox()
		for _, child := range children {
			box.Append(schwifty.ResolveWidget(child))
		}
		return box
	})
}
