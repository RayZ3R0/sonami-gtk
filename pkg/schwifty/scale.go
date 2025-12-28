package schwifty

import (
	"codeberg.org/dergs/tidalwave/internal/g"
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
		label := f()

		var callbackId string
		label.ConnectRealize(g.Ptr(func(a gtk.Widget) {
			callbackId = state.AddCallback(func(newValue float64) {
				gtk.ScaleNewFromInternalPtr(a.GoPointer()).SetValue(newValue)
			})
		}))
		label.ConnectUnrealize(g.Ptr(func(gtk.Widget) {
			state.RemoveCallback(callbackId)
		}))

		return label
	}
}

func (s Scale) ConnectChangeValue(callback func(gtk.Range, gtk.ScrollType, float64) bool) Scale {
	return func() *gtk.Scale {
		scale := s()
		scale.ConnectChangeValue(&callback)
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
