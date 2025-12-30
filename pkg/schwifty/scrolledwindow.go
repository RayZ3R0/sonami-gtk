package schwifty

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

//go:generate go run codeberg.org/dergs/tidalwave/pkg/schwifty/gen ScrolledWindow *gtk.ScrolledWindow

func (f ScrolledWindow) BindChild(state *state.State[any]) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		var callbackId string
		return f.ConnectRealize(func(w gtk.Widget) {
			callbackId = state.AddCallback(func(newValue any) {
				widget := ResolveWidget(newValue)
				widget.Ref()
				OnMainThreadOnce(func(u uintptr) {
					gtk.ScrolledWindowNewFromInternalPtr(u).SetChild(widget)
					widget.Unref()
				}, w.GoPointer())
			})
		}).ConnectUnrealize(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

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

func (f ScrolledWindow) PropagateNaturalHeight(propagate bool) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		scrolledWindow := f()
		scrolledWindow.SetPropagateNaturalHeight(propagate)
		return scrolledWindow
	}
}

func (f ScrolledWindow) PropagateNaturalWidth(propagate bool) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		scrolledWindow := f()
		scrolledWindow.SetPropagateNaturalWidth(propagate)
		return scrolledWindow
	}
}
