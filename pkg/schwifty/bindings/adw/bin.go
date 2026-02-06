package adw

import (
	gtkbindings "codeberg.org/dergs/tonearm/pkg/schwifty/bindings/gtk"
	"codeberg.org/dergs/tonearm/pkg/schwifty/callback"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	"codeberg.org/dergs/tonearm/pkg/schwifty/tracking"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

//go:generate go run codeberg.org/dergs/tonearm/pkg/schwifty/gen Bin *adw.Bin adw

func (f Bin) BindChild(state *state.State[any]) Bin {
	return func() *adw.Bin {
		var callbackId string
		var ref *tracking.WeakRef
		return f.ConnectConstruct(func(w *adw.Bin) {
			ref = tracking.NewWeakRef(w)
			callbackId = state.AddCallback(func(newValue any) {
				widget := gtkbindings.ResolveWidget(newValue)
				if widget == nil {
					callback.OnMainThreadOncePure(func() {
						if obj := ref.Get(); obj != nil {
							defer obj.Unref()
							adw.BinNewFromInternalPtr(obj.Ptr).SetChild(nil)
						}
					})
				} else {
					widget.Ref()
					callback.OnMainThreadOncePure(func() {
						defer widget.Unref()
						if obj := ref.Get(); obj != nil {
							defer obj.Unref()
							adw.BinNewFromInternalPtr(obj.Ptr).SetChild(widget)
						}
					})

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
