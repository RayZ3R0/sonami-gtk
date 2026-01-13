package schwifty

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty/callback"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

//go:generate go run codeberg.org/dergs/tonearm/pkg/schwifty/gen ScrolledWindow *gtk.ScrolledWindow

func (f ScrolledWindow) BindChild(state *state.State[any]) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.ScrolledWindow) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue any) {
				widget := ResolveWidget(newValue)
				if widget == nil {
					OnMainThreadOnce(func(u uintptr) {
						gtk.ScrolledWindowNewFromInternalPtr(u).SetChild(nil)
					}, widgetPtr)
				} else {
					widget.Ref()
					OnMainThreadOnce(func(u uintptr) {
						gtk.ScrolledWindowNewFromInternalPtr(u).SetChild(widget)
						widget.Unref()
					}, widgetPtr)
				}
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f ScrolledWindow) ConnectEdgeReached(cb func(gtk.ScrolledWindow, gtk.PositionType)) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		scrolledWindow := f()
		callback.HandleCallback(scrolledWindow.Object, "edge-reached", cb)
		return scrolledWindow
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
