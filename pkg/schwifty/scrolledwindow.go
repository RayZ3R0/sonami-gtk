package schwifty

import "github.com/jwijenbergh/puregotk/v4/gtk"

//go:generate go run codeberg.org/dergs/tidalwave/pkg/schwifty/gen ScrolledWindow *gtk.ScrolledWindow

func (f ScrolledWindow) Child(widget any) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		scrolledWindow := f()
		scrolledWindow.SetChild(ResolveWidget(widget))
		return scrolledWindow
	}
}

func (f ScrolledWindow) Policy(hPolicy, vPolicy gtk.PolicyType) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		scrolledWindow := f()
		scrolledWindow.SetPolicy(hPolicy, vPolicy)
		return scrolledWindow
	}
}
