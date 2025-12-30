package schwifty

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

//go:generate go run codeberg.org/dergs/tidalwave/pkg/schwifty/gen CenterBox *gtk.CenterBox

func (f CenterBox) BindCenterWidget(state *state.State[any]) CenterBox {
	return func() *gtk.CenterBox {
		var callbackId string
		return f.ConnectRealize(func(w gtk.Widget) {
			callbackId = state.AddCallback(func(newValue any) {
				widget := ResolveWidget(newValue)
				if widget == nil {
					OnMainThreadOnce(func(u uintptr) {
						gtk.CenterBoxNewFromInternalPtr(u).SetCenterWidget(nil)
					}, w.GoPointer())
				} else {
					widget.Ref()
					OnMainThreadOnce(func(u uintptr) {
						gtk.CenterBoxNewFromInternalPtr(u).SetCenterWidget(widget)
						widget.Unref()
					}, w.GoPointer())
				}
			})
		}).ConnectUnrealize(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f CenterBox) CenterWidget(child any) CenterBox {
	return func() *gtk.CenterBox {
		centerBox := f()
		centerBox.SetCenterWidget(ResolveWidget(child))
		return centerBox
	}
}

func (f CenterBox) EndWidget(child any) CenterBox {
	return func() *gtk.CenterBox {
		centerBox := f()
		centerBox.SetEndWidget(ResolveWidget(child))
		return centerBox
	}
}

func (f CenterBox) StartWidget(child any) CenterBox {
	return func() *gtk.CenterBox {
		centerBox := f()
		centerBox.SetStartWidget(ResolveWidget(child))
		return centerBox
	}
}
