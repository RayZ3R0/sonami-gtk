package syntax

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"github.com/jwijenbergh/puregotk/v4/adw"
)

func SwitchRow() schwifty.SwitchRow {
	return managed("SwitchRow", func() *adw.SwitchRow {
		return adw.NewSwitchRow()
	})
}
