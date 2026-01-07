package syntax

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"github.com/jwijenbergh/puregotk/v4/adw"
)

func PreferencesGroup(children ...any) schwifty.PreferencesGroup {
	return managed("PreferencesGroup", func() *adw.PreferencesGroup {
		group := adw.NewPreferencesGroup()
		for _, child := range children {
			group.Add(schwifty.ResolveWidget(child))
		}
		return group
	})
}
