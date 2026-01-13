package schwifty

import (
	"github.com/jwijenbergh/puregotk/v4/adw"
)

//go:generate go run codeberg.org/dergs/tonearm/pkg/schwifty/gen PreferencesPage *adw.PreferencesPage

func (f PreferencesPage) Add(child any) PreferencesPage {
	return func() *adw.PreferencesPage {
		page := f()
		page.Add(ResolveTo[*adw.PreferencesGroup, PreferencesGroup](child))
		return page
	}
}

func (f PreferencesPage) IconName(iconName string) PreferencesPage {
	return func() *adw.PreferencesPage {
		page := f()
		page.SetIconName(iconName)
		return page
	}
}

func (f PreferencesPage) Title(title string) PreferencesPage {
	return func() *adw.PreferencesPage {
		page := f()
		page.SetTitle(title)
		return page
	}
}
