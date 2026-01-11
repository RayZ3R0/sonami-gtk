package syntax

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"github.com/jwijenbergh/puregotk/v4/adw"
)

func PasswordEntryRow() schwifty.PasswordEntryRow {
	return managed("PasswordEntryRow", func() *adw.PasswordEntryRow {
		return adw.NewPasswordEntryRow()
	})
}
