package gtk

import (
	"github.com/jwijenbergh/puregotk/v4/gdk"
	"github.com/jwijenbergh/puregotk/v4/gobject"
)

type Paintable interface {
	gdk.Paintable
	Ref() *gobject.Object
	Unref()
}
