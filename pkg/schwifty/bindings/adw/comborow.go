package adw

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty/callback"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gio"
	"github.com/jwijenbergh/puregotk/v4/gobject"
)

var (
	ComboRowSelectionChangedCallback = func(obj gobject.Object, _ uintptr) {
		callback.CallbackHandler[any](obj, "notify", adw.ComboRowNewFromInternalPtr(obj.Ptr))
	}
)

//go:generate go run codeberg.org/dergs/tonearm/pkg/schwifty/gen ComboRow *adw.ComboRow adw

func (f ComboRow) ConnectSelectionChanged(cb func(uint)) ComboRow {
	return func() *adw.ComboRow {
		row := f()
		var selection uint
		callback.HandleCallback(row.Object, "notify", func(comboRow *adw.ComboRow) {
			newValue := comboRow.GetSelected()
			if newValue != selection {
				selection = newValue
				cb(selection)
			}
		})
		return row
	}
}

func (f ComboRow) Model(model gio.ListModel) ComboRow {
	return func() *adw.ComboRow {
		row := f()
		row.SetModel(model)
		return row
	}
}

func (f ComboRow) Selected(index uint) ComboRow {
	return func() *adw.ComboRow {
		row := f()
		row.SetSelected(index)
		return row
	}
}

func (f ComboRow) Subtitle(subtitle string) ComboRow {
	return func() *adw.ComboRow {
		row := f()
		row.SetSubtitle(subtitle)
		return row
	}
}

func (f ComboRow) Title(title string) ComboRow {
	return func() *adw.ComboRow {
		row := f()
		row.SetTitle(title)
		return row
	}
}
