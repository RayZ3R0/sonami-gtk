package schwifty

import "github.com/jwijenbergh/puregotk/v4/adw"

//go:generate go run codeberg.org/dergs/tonearm/pkg/schwifty/gen SwitchRow *adw.SwitchRow

func (f SwitchRow) Subtitle(subtitle string) SwitchRow {
	return func() *adw.SwitchRow {
		switchRow := f()
		switchRow.SetSubtitle(subtitle)
		return switchRow
	}
}

func (f SwitchRow) Title(title string) SwitchRow {
	return func() *adw.SwitchRow {
		switchRow := f()
		switchRow.SetTitle(title)
		return switchRow
	}
}
