package syntax

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"github.com/jwijenbergh/puregotk/v4/adw"
)

func HeaderBar() schwifty.HeaderBar {
	return managed("HeaderBar", func() *adw.HeaderBar {
		return adw.NewHeaderBar()
	})
}
