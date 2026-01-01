package schwifty

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

//go:generate go run codeberg.org/dergs/tidalwave/pkg/schwifty/gen MenuButton *gtk.MenuButton

func (f MenuButton) BindIconName(state *state.State[string]) MenuButton {
	return func() *gtk.MenuButton {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.MenuButton) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue string) {
				OnMainThreadOnce(func(u uintptr) {
					gtk.MenuButtonNewFromInternalPtr(u).SetIconName(newValue)
				}, widgetPtr)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
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
