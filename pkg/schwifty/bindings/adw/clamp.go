package adw

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty/bindings/gtk"
	"github.com/jwijenbergh/puregotk/v4/adw"
)

//go:generate go run codeberg.org/dergs/tonearm/pkg/schwifty/gen Clamp *adw.Clamp adw

func (f Clamp) Child(widget any) Clamp {
	return func() *adw.Clamp {
		clamp := f()
		clamp.SetChild(gtk.ResolveWidget(widget))
		return clamp
	}
}

func (f Clamp) MaximumSize(size int) Clamp {
	return func() *adw.Clamp {
		clamp := f()
		clamp.SetMaximumSize(size)
		return clamp
	}
}
