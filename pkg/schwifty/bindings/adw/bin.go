package adw

import (
	gtkbindings "codeberg.org/dergs/tonearm/pkg/schwifty/bindings/gtk"
	"github.com/jwijenbergh/puregotk/v4/adw"
)

//go:generate go run codeberg.org/dergs/tonearm/pkg/schwifty/gen Bin *adw.Bin adw

func (f Bin) Child(widget any) Bin {
	return func() *adw.Bin {
		bin := f()
		bin.SetChild(gtkbindings.ResolveWidget(widget))
		return bin
	}
}
