package syntax

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"github.com/jwijenbergh/puregotk/v4/adw"
)

func PreferencesPage(groups ...any) schwifty.PreferencesPage {
	return managed("PreferencesPage", func() *adw.PreferencesPage {
		page := adw.NewPreferencesPage()
		for _, group := range groups {
			page.Add(schwifty.ResolveTo[*adw.PreferencesGroup, schwifty.PreferencesGroup](group))
		}
		return page
	})
}
