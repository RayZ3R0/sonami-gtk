package schwifty

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty/callback"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

//go:generate go run codeberg.org/dergs/tidalwave/pkg/schwifty/gen Scale *gtk.Scale

func (s Scale) Value(value float64) Scale {
	return func() *gtk.Scale {
		scale := s()
		scale.SetValue(value)
		return scale
	}
}

func (f Scale) BindValue(state *state.State[float64]) Scale {
	return func() *gtk.Scale {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.Scale) {
			callbackId = state.AddCallback(func(newValue float64) {
				gtk.ScaleNewFromInternalPtr(w.GoPointer()).SetValue(newValue)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (s Scale) ConnectChangeValue(cb func(gtk.Range, gtk.ScrollType, float64) bool) Scale {
	return func() *gtk.Scale {
		scale := s()
		callback.HandleCallback(scale.Widget, "change-value", cb)
		return scale
	}
}

func (s Scale) Inverted(invert bool) Scale {
	return func() *gtk.Scale {
		scale := s()
		scale.SetInverted(invert)
		return scale
	}
}

func (s Scale) Range(min, max float64) Scale {
	return func() *gtk.Scale {
		scale := s()
		scale.SetRange(min, max)
		return scale
	}
}
