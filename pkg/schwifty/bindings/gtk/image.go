package gtk

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty/callback"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	"codeberg.org/dergs/tonearm/pkg/schwifty/tracking"
	"github.com/jwijenbergh/puregotk/v4/gdk"
	"github.com/jwijenbergh/puregotk/v4/gdkpixbuf"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

//go:generate go run codeberg.org/dergs/tonearm/pkg/schwifty/gen Image *gtk.Image gtk

func (f Image) BindPaintable(state *state.State[Paintable]) Image {
	return func() *gtk.Image {
		var callbackId string
		var ref *tracking.WeakRef
		return f.ConnectConstruct(func(w *gtk.Image) {
			ref = tracking.NewWeakRef(w)
			callbackId = state.AddCallback(func(newValue Paintable) {
				newValue.Ref()
				callback.OnMainThreadOncePure(func() {
					defer newValue.Unref()
					if obj := ref.Get(); obj != nil {
						defer obj.Unref()
						gtk.ImageNewFromInternalPtr(obj.Ptr).SetFromPaintable(newValue)
					}
				})
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f Image) BindPixbuf(state *state.State[*gdkpixbuf.Pixbuf]) Image {
	return func() *gtk.Image {
		var callbackId string
		var ref *tracking.WeakRef
		return f.ConnectConstruct(func(w *gtk.Image) {
			ref = tracking.NewWeakRef(w)
			callbackId = state.AddCallback(func(newValue *gdkpixbuf.Pixbuf) {
				newValue.Ref()
				callback.OnMainThreadOncePure(func() {
					defer newValue.Unref()
					if obj := ref.Get(); obj != nil {
						defer obj.Unref()
						gtk.ImageNewFromInternalPtr(obj.Ptr).SetFromPixbuf(newValue)
					}
				})
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f Image) FromIconName(iconName string) Image {
	return func() *gtk.Image {
		image := f()
		image.SetFromIconName(iconName)
		return image
	}
}

func (f Image) FromPaintable(paintable gdk.Paintable) Image {
	return func() *gtk.Image {
		image := f()
		image.SetFromPaintable(paintable)
		return image
	}
}

func (f Image) FromPixbuf(pixbuf *gdkpixbuf.Pixbuf) Image {
	return func() *gtk.Image {
		image := f()
		image.SetFromPixbuf(pixbuf)
		return image
	}
}

func (f Image) FromResource(resource string) Image {
	return func() *gtk.Image {
		image := f()
		image.SetFromResource(resource)
		return image
	}
}

func (f Image) PixelSize(size int) Image {
	return func() *gtk.Image {
		image := f()
		image.SetPixelSize(size)
		return image
	}
}
