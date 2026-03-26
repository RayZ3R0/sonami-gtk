package adw

import (
	"codeberg.org/puregotk/puregotk/v4/adw"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/bindings"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/bindings/gtk"
)

//go:generate go run github.com/RayZ3R0/sonami-gtk/pkg/schwifty/gen PreferencesDialog *adw.PreferencesDialog adw

func (f PreferencesDialog) Add(child any) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		dialog := f()
		dialog.Add(bindings.ResolveTo[*adw.PreferencesPage, PreferencesPage](child))
		return dialog
	}
}

func (f PreferencesDialog) Present(parent any) {
	resolvedParent := gtk.ResolveWidget(parent)
	if resolvedParent != nil {
		dialog := f()
		dialog.Present(resolvedParent)
	}
}
