package schwifty

import "github.com/jwijenbergh/puregotk/v4/adw"

//go:generate go run codeberg.org/dergs/tidalwave/pkg/schwifty/gen SpinRow *adw.SpinRow

func (f SpinRow) Title(title string) SpinRow {
	return func() *adw.SpinRow {
		row := f()
		row.SetTitle(title)
		return row
	}
}

func (f SpinRow) Subtitle(subtitle string) SpinRow {
	return func() *adw.SpinRow {
		row := f()
		row.SetSubtitle(subtitle)
		return row
	}
}
