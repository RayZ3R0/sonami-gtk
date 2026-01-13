package syntax

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"github.com/jwijenbergh/puregotk/v4/adw"
)

func PreferencesDialog(pages ...any) schwifty.PreferencesDialog {
	return managed("PreferencesDialog", func() *adw.PreferencesDialog {
		dialog := adw.NewPreferencesDialog()
		for _, page := range pages {
			dialog.Add(schwifty.ResolveTo[*adw.PreferencesPage, schwifty.PreferencesPage](page))
		}
		return dialog
	})
}
