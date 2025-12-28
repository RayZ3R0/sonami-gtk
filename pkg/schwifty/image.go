package schwifty

import (
	"codeberg.org/dergs/tidalwave/internal/g"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	"github.com/jwijenbergh/puregotk/v4/gdkpixbuf"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

//go:generate go run codeberg.org/dergs/tidalwave/pkg/schwifty/gen Image *gtk.Image

func (f Image) BindPixbuf(state *state.State[*gdkpixbuf.Pixbuf]) Image {
	return func() *gtk.Image {
		image := f()

		var callbackId string
		image.ConnectRealize(g.Ptr(func(a gtk.Widget) {
			callbackId = state.AddCallback(func(newValue *gdkpixbuf.Pixbuf) {
				gtk.ImageNewFromInternalPtr(a.GoPointer()).SetFromPixbuf(newValue)
			})
		}))
		image.ConnectUnrealize(g.Ptr(func(gtk.Widget) {
			state.RemoveCallback(callbackId)
		}))

		return image
	}
}

func (f Image) FromIconName(iconName string) Image {
	return func() *gtk.Image {
		image := f()
		image.SetFromIconName(iconName)
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
