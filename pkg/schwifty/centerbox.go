package schwifty

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

//go:generate go run codeberg.org/dergs/tonearm/pkg/schwifty/gen CenterBox *gtk.CenterBox

func (f CenterBox) BindCenterWidget(state *state.State[any]) CenterBox {
	return func() *gtk.CenterBox {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.CenterBox) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue any) {
				widget := ResolveWidget(newValue)
				if widget == nil {
					OnMainThreadOnce(func(u uintptr) {
						gtk.CenterBoxNewFromInternalPtr(u).SetCenterWidget(nil)
					}, widgetPtr)
				} else {
					widget.Ref()
					OnMainThreadOnce(func(u uintptr) {
						gtk.CenterBoxNewFromInternalPtr(u).SetCenterWidget(widget)
						widget.Unref()
					}, widgetPtr)
				}
			})
		}).ConnectDestroy(func(w gtk.Widget) {
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
