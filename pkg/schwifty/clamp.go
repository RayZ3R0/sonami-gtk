package schwifty

import "github.com/jwijenbergh/puregotk/v4/adw"

//go:generate go run codeberg.org/dergs/tidalwave/pkg/schwifty/gen Clamp *adw.Clamp

func (f Clamp) Child(widget any) Clamp {
	return func() *adw.Clamp {
		clamp := f()
		clamp.SetChild(ResolveWidget(widget))
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
