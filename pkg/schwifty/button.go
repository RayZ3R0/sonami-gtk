package schwifty

import "github.com/jwijenbergh/puregotk/v4/gtk"

//go:generate go run codeberg.org/dergs/tidalwave/pkg/schwifty/gen Button *gtk.Button

func (f Button) Child(widget any) Button {
	return func() *gtk.Button {
		clamp := f()
		clamp.SetChild(ResolveWidget(widget))
		return clamp
	}
}

func (f Button) ConnectClicked(cb func(gtk.Button)) Button {
	return func() *gtk.Button {
		clamp := f()
		clamp.ConnectClicked(&cb)
		return clamp
	}
}
