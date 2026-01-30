package adw

import (
	gtkbindings "codeberg.org/dergs/tonearm/pkg/schwifty/bindings/gtk"
	"codeberg.org/dergs/tonearm/pkg/schwifty/callback"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

//go:generate go run codeberg.org/dergs/tonearm/pkg/schwifty/gen Bin *adw.Bin adw

func (f Bin) BindChild(state *state.State[any]) Bin {
	return func() *adw.Bin {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.Bin) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue any) {
				widget := gtkbindings.ResolveWidget(newValue)
				if widget == nil {
					callback.OnMainThreadOnce(func(u uintptr) {
						adw.BinNewFromInternalPtr(u).SetChild(nil)
					}, widgetPtr)
				} else {
					widget.Ref()
					callback.OnMainThreadOnce(func(u uintptr) {
						adw.BinNewFromInternalPtr(u).SetChild(widget)
						widget.Unref()
					}, widgetPtr)
				}
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f Bin) Child(widget any) Bin {
	return func() *adw.Bin {
		bin := f()
		bin.SetChild(gtkbindings.ResolveWidget(widget))
		return bin
	}
}
