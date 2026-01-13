package syntax

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"github.com/jwijenbergh/puregotk/v4/adw"
)

func StatusPage() schwifty.StatusPage {
	return managed("StatusPage", func() *adw.StatusPage {
		return adw.NewStatusPage()
	})
}
