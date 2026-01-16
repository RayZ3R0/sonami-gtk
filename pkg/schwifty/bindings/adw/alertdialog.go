package adw

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty/bindings/gtk"
	"codeberg.org/dergs/tonearm/pkg/schwifty/callback"
	"github.com/jwijenbergh/puregotk/v4/adw"
)

//go:generate go run codeberg.org/dergs/tonearm/pkg/schwifty/gen AlertDialog *adw.AlertDialog adw

func (f AlertDialog) CanClose(canClose bool) AlertDialog {
	return func() *adw.AlertDialog {
		dialog := f()
		dialog.SetCanClose(canClose)
		return dialog
	}
}

func (f AlertDialog) ConnectCloseAttempt(cb func(adw.Dialog)) AlertDialog {
	return func() *adw.AlertDialog {
		dialog := f()
		callback.HandleCallback(dialog.Object, "close-attempt", cb)
		return dialog
	}
}

func (f AlertDialog) ExtraChild(widget any) AlertDialog {
	return func() *adw.AlertDialog {
		dialog := f()
		dialog.SetExtraChild(gtk.ResolveWidget(widget))
		return dialog
	}
}
