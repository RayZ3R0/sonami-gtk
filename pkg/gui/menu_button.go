package gui

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type MenuButtonImpl struct {
	*WidgetImpl[*MenuButtonImpl]
	menuButton *gtk.MenuButton
}

func MenuButton(popover *PopoverImpl) *MenuButtonImpl {
	button := gtk.NewMenuButton()
	button.SetPopover(popover)

	impl := &MenuButtonImpl{nil, button}
	impl.WidgetImpl = &WidgetImpl[*MenuButtonImpl]{button, button.Widget, impl}
	return impl
}

func (b *MenuButtonImpl) SetIcon(iconName string) *MenuButtonImpl {
	b.menuButton.SetIconName(iconName)
	return b
}
