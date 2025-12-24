package gui

import "github.com/diamondburned/gotk4/pkg/gtk/v4"

type AspectFrameFunc func(child gtk.Widgetter) *AspectFrameImpl

type AspectFrameImpl struct {
	*WidgetImpl[*AspectFrameImpl]
	frame *gtk.AspectFrame
}

var AspectFrame = func(child gtk.Widgetter) *AspectFrameImpl {
	frame := gtk.NewAspectFrame(0.5, 0.5, 1.0, false)
	frame.SetChild(child)
	frame.SetOverflow(gtk.OverflowHidden)
	impl := &AspectFrameImpl{nil, frame}
	impl.WidgetImpl = &WidgetImpl[*AspectFrameImpl]{frame, frame.Widget, impl}
	return impl
}
