package schwifty

import (
	"github.com/jwijenbergh/puregotk/v4/gdkpixbuf"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

//go:generate go run codeberg.org/dergs/tidalwave/pkg/schwifty/gen Image *gtk.Image

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
