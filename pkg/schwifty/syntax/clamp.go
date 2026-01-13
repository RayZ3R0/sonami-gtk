package syntax

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"github.com/jwijenbergh/puregotk/v4/adw"
)

func Clamp() schwifty.Clamp {
	return managed("Clamp", func() *adw.Clamp {
		return adw.NewClamp()
	})
}
