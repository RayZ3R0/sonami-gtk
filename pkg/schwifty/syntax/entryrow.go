package syntax

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"github.com/jwijenbergh/puregotk/v4/adw"
)

func EntryRow() schwifty.EntryRow {
	return managed("EntryRow", func() *adw.EntryRow {
		return adw.NewEntryRow()
	})
}
