package syntax

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func SpinRow(adjustment *gtk.Adjustment, climbRate float64, digits uint) schwifty.SpinRow {
	return managed("SpinRow", func() *adw.SpinRow {
		return adw.NewSpinRow(adjustment, climbRate, digits)
	})
}
