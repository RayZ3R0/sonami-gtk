package gtk

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty/bindings"
	"codeberg.org/dergs/tonearm/pkg/schwifty/callback"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	"codeberg.org/dergs/tonearm/pkg/schwifty/tracking"
	"github.com/jwijenbergh/puregotk/v4/gio"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

//go:generate go run codeberg.org/dergs/tonearm/pkg/schwifty/gen MenuButton *gtk.MenuButton gtk

func (f MenuButton) BindIconName(state *state.State[string]) MenuButton {
	return func() *gtk.MenuButton {
		var callbackId string
		var ref *tracking.WeakRef
		return f.ConnectConstruct(func(w *gtk.MenuButton) {
			ref = tracking.NewWeakRef(w)
			callbackId = state.AddCallback(func(newValue string) {
				callback.OnMainThreadOncePure(func() {
					if obj := ref.Get(); obj != nil {
						defer obj.Unref()
						gtk.MenuButtonNewFromInternalPtr(obj.Ptr).SetIconName(newValue)
					}
				})
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

func (f MenuButton) MenuModel(model *gio.MenuModel) MenuButton {
	return func() *gtk.MenuButton {
		button := f()
		button.SetMenuModel(model)
		return button
	}
}

func (f MenuButton) Popover(widget any) MenuButton {
	return func() *gtk.MenuButton {
		button := f()
		button.SetPopover(bindings.ResolveTo[*gtk.Popover, Popover](widget))
		return button
	}
}

func (f MenuButton) TooltipText(tooltip string) MenuButton {
	return func() *gtk.MenuButton {
		button := f()
		button.SetTooltipText(tooltip)
		return button
	}
}
