package gtk

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty/callback"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	"codeberg.org/dergs/tonearm/pkg/schwifty/tracking"
	"github.com/jwijenbergh/puregotk/v4/gdk"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

//go:generate go run codeberg.org/dergs/tonearm/pkg/schwifty/gen Picture *gtk.Picture gtk

func (f Picture) BindPaintable(state *state.State[Paintable]) Picture {
	return func() *gtk.Picture {
		var callbackId string
		var ref *tracking.WeakRef
		return f.ConnectConstruct(func(w *gtk.Picture) {
			ref = tracking.NewWeakRef(w)
			callbackId = state.AddCallback(func(newValue Paintable) {
				newValue.Ref()
				callback.OnMainThreadOncePure(func() {
					defer newValue.Unref()
					if obj := ref.Get(); obj != nil {
						defer obj.Unref()
						gtk.PictureNewFromInternalPtr(obj.Ptr).SetPaintable(newValue)
					}
				})
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f Picture) ContentFit(fit gtk.ContentFit) Picture {
	return func() *gtk.Picture {
		picture := f()
		picture.SetContentFit(fit)
		return picture
	}
}

func (f Picture) FromPaintable(paintable gdk.Paintable) Picture {
	return func() *gtk.Picture {
		picture := f()
		picture.SetPaintable(paintable)
		return picture
	}
}
