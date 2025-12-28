package syntax

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func Box(orientation gtk.Orientation, children ...any) schwifty.Box {
	return managed("Box", func() *gtk.Box {
		box := gtk.NewBox(orientation, 0)
		for _, child := range children {
			box.Append(schwifty.ResolveWidget(child))
		}
		return box
	})
}

func HStack(children ...any) schwifty.Box {
	return Box(gtk.OrientationHorizontalValue, children...)
}

func Spacer() schwifty.Box {
	return HStack().VExpand(true).HExpand(true)
}

func VStack(children ...any) schwifty.Box {
	return Box(gtk.OrientationVerticalValue, children...)
}
