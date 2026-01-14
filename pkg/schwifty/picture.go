package schwifty

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	"github.com/jwijenbergh/puregotk/v4/gdk"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

//go:generate go run codeberg.org/dergs/tonearm/pkg/schwifty/gen Picture *gtk.Picture

func (f Picture) BindPaintable(state *state.State[Paintable]) Picture {
	return func() *gtk.Picture {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.Picture) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue Paintable) {
				newValue.Ref()
				OnMainThreadOnce(func(u uintptr) {
					gtk.PictureNewFromInternalPtr(u).SetPaintable(newValue)
					newValue.Unref()
				}, widgetPtr)
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
