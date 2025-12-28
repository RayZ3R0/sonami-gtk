package schwifty

import (
	"codeberg.org/dergs/tidalwave/internal/g"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func (f TEMPLATE_TYPE) BindVisible(state *state.State[bool]) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		widget := f()

		var callbackId string
		widget.ConnectRealize(g.Ptr(func(a gtk.Widget) {
			callbackId = state.AddCallback(func(newValue bool) {
				a.SetVisible(newValue)
			})
		}))
		widget.ConnectUnrealize(g.Ptr(func(gtk.Widget) {
			state.RemoveCallback(callbackId)
		}))

		return widget
	}
}

func (f TEMPLATE_TYPE) BindSensitive(state *state.State[bool]) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		widget := f()

		var callbackId string
		widget.ConnectRealize(g.Ptr(func(a gtk.Widget) {
			callbackId = state.AddCallback(func(newValue bool) {
				a.SetSensitive(newValue)
			})
		}))
		widget.ConnectUnrealize(g.Ptr(func(gtk.Widget) {
			state.RemoveCallback(callbackId)
		}))

		return widget
	}
}
