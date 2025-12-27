package gui

import "github.com/diamondburned/gotk4/pkg/gtk/v4"

type PopoverImpl struct {
	*WidgetImpl[*PopoverImpl]
	popover *gtk.Popover
}

func Popover(child gtk.Widgetter) *PopoverImpl {
	popover := gtk.NewPopover()
	popover.SetChild(child)

	impl := &PopoverImpl{nil, popover}
	impl.WidgetImpl = &WidgetImpl[*PopoverImpl]{popover, popover.Widget, impl}
	return impl
}
