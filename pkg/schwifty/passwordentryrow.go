package schwifty

import "github.com/jwijenbergh/puregotk/v4/adw"

//go:generate go run codeberg.org/dergs/tidalwave/pkg/schwifty/gen PasswordEntryRow *adw.PasswordEntryRow

func (f PasswordEntryRow) Title(title string) PasswordEntryRow {
	return func() *adw.PasswordEntryRow {
		row := f()
		row.SetTitle(title)
		return row
	}
}
