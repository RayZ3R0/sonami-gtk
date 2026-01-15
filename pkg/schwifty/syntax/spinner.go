package syntax

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func Spinner() schwifty.Spinner {
	return managed("Spinner", func() *gtk.Spinner {
		spinner := gtk.NewSpinner()
		spinner.Start()
		return spinner
	})
}
