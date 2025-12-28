package schwifty

import (
	"codeberg.org/dergs/tidalwave/internal/g"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

//go:generate go run codeberg.org/dergs/tidalwave/pkg/schwifty/gen MenuButton *gtk.MenuButton

func (f MenuButton) BindIconName(state *state.State[string]) MenuButton {
	return func() *gtk.MenuButton {
		label := f()

		var callbackId string
		label.ConnectRealize(g.Ptr(func(a gtk.Widget) {
			callbackId = state.AddCallback(func(newValue string) {
				gtk.MenuButtonNewFromInternalPtr(a.GoPointer()).SetIconName(newValue)
			})
		}))
		label.ConnectUnrealize(g.Ptr(func(gtk.Widget) {
			state.RemoveCallback(callbackId)
		}))

		return label
	}
}

func (f MenuButton) Child(widget any) MenuButton {
	return func() *gtk.MenuButton {
		button := f()
		button.SetChild(ResolveWidget(widget))
		return button
	}
}

func (f MenuButton) IconName(iconName string) MenuButton {
	return func() *gtk.MenuButton {
		button := f()
		button.SetIconName(iconName)
		return button
	}
}

func (f MenuButton) Popover(widget any) MenuButton {
	return func() *gtk.MenuButton {
		button := f()
		button.SetPopover(ResolvePopover(widget))
		return button
	}
}
