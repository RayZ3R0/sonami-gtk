package schwifty

import (
	"github.com/jwijenbergh/puregotk/v4/adw"
)

//go:generate go run codeberg.org/dergs/tidalwave/pkg/schwifty/gen PreferencesDialog *adw.PreferencesDialog

func (f PreferencesDialog) Add(child any) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		dialog := f()
		dialog.Add(ResolveTo[*adw.PreferencesPage, PreferencesPage](child))
		return dialog
	}
}

func (f PreferencesDialog) Present(parent any) {
	resolvedParent := ResolveWidget(parent)
	if resolvedParent != nil {
		dialog := f()
		dialog.Present(resolvedParent)
	}
}
