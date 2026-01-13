package schwifty

import (
	"github.com/jwijenbergh/puregotk/v4/adw"
)

//go:generate go run codeberg.org/dergs/tonearm/pkg/schwifty/gen PreferencesGroup *adw.PreferencesGroup

func (f PreferencesGroup) Add(child any) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		group := f()
		group.Add(ResolveWidget(child))
		return group
	}
}

func (f PreferencesGroup) Description(description string) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		group := f()
		group.SetDescription(description)
		return group
	}
}

func (f PreferencesGroup) Title(title string) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		group := f()
		group.SetTitle(title)
		return group
	}
}
