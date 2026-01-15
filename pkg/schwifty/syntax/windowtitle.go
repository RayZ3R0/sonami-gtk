package syntax

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"github.com/jwijenbergh/puregotk/v4/adw"
)

func WindowTitle(title string, subtitle string) schwifty.WindowTitle {
	return managed("WindowTitle", func() *adw.WindowTitle {
		return adw.NewWindowTitle(title, subtitle)
	})
}
