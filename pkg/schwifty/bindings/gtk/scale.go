package gtk

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty/callback"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	"codeberg.org/dergs/tonearm/pkg/schwifty/tracking"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

//go:generate go run codeberg.org/dergs/tonearm/pkg/schwifty/gen Scale *gtk.Scale gtk

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
		var ref *tracking.WeakRef
		return f.ConnectConstruct(func(w *gtk.Scale) {
			ref = tracking.NewWeakRef(w)
			callbackId = state.AddCallback(func(newValue float64) {
				callback.OnMainThreadOncePure(func() {
					if obj := ref.Get(); obj != nil {
						defer obj.Unref()
						gtk.ScaleNewFromInternalPtr(obj.GoPointer()).SetValue(newValue)
					}
				})
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (s Scale) ConnectChangeValue(cb func(gtk.Range, gtk.ScrollType, float64) bool) Scale {
	return func() *gtk.Scale {
		scale := s()
		callback.HandleCallback(scale.Object, "change-value", cb)
		return scale
	}
}

func (s Scale) Increments(stepSize float64, pageSize float64) Scale {
	return func() *gtk.Scale {
		scale := s()
		scale.SetIncrements(stepSize, pageSize)
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
