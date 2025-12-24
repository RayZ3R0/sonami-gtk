package gui

import (
	"fmt"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotkit/gtkutil/cssutil"
)

type BoxFunc func(children ...gtk.Widgetter) *BoxImpl

type BoxImpl struct {
	*WidgetImpl[*BoxImpl]
	box *gtk.Box
}

var Box = func(children ...gtk.Widgetter) *BoxImpl {
	box := gtk.NewBox(gtk.OrientationVertical, 0)
	for _, child := range children {
		box.Append(child)
	}

	impl := &BoxImpl{nil, box}
	impl.WidgetImpl = &WidgetImpl[*BoxImpl]{box, box.Widget, impl}
	return impl
}

var HStack = func(children ...gtk.Widgetter) *BoxImpl {
	return Box(children...).Orientation(gtk.OrientationHorizontal)
}

var VStack = func(children ...gtk.Widgetter) *BoxImpl {
	return Box(children...).Orientation(gtk.OrientationVertical)
}

var Spacer = func() *BoxImpl {
	return Box().VExpand(true).HExpand(true)
}

func (b *BoxImpl) Append(child gtk.Widgetter) *BoxImpl {
	b.box.Append(child)
	return b
}

func (b *BoxImpl) GTKWidget() *gtk.Box {
	return b.box
}

func (b *BoxImpl) Orientation(orientation gtk.Orientation) *BoxImpl {
	b.box.SetOrientation(orientation)
	return b
}

func (b *BoxImpl) Spacing(spacing int) *BoxImpl {
	b.box.SetSpacing(spacing)
	return b
}

func (b *BoxImpl) MinHeight(minHeightPx int) *BoxImpl {
	cssutil.Apply(b, fmt.Sprintf("%s { min-height: %dpx; }", b.widget.CSSName(), minHeightPx))
	return b
}

func (b *BoxImpl) MinWidth(minWidthPx int) *BoxImpl {
	cssutil.Apply(b, fmt.Sprintf("%s { min-width: %dpx; }", b.widget.CSSName(), minWidthPx))
	return b
}
