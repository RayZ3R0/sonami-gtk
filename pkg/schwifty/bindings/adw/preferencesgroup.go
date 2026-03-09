package adw

import (
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/bindings/gtk"
	"codeberg.org/puregotk/puregotk/v4/adw"
)

//go:generate go run github.com/RayZ3R0/sonami-gtk/pkg/schwifty/gen PreferencesGroup *adw.PreferencesGroup adw

func (f PreferencesGroup) Add(child any) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		group := f()
		group.Add(gtk.ResolveWidget(child))
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
