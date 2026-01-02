package schwifty

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty/callback"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	"github.com/jwijenbergh/puregotk/v4/glib"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

//go:generate go run codeberg.org/dergs/tidalwave/pkg/schwifty/gen Button *gtk.Button

func (f Button) ActionName(actionName string) Button {
	return func() *gtk.Button {
		button := f()
		button.SetActionName(actionName)
		return button
	}
}

func (f Button) ActionTargetValue(targetValue *glib.Variant) Button {
	return func() *gtk.Button {
		button := f()
		button.SetActionTargetValue(targetValue)
		return button
	}
}

func (f Button) BindIconName(state *state.State[string]) Button {
	return func() *gtk.Button {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.Button) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue string) {
				callback.OnMainThreadOnce(func(u uintptr) {
					gtk.ButtonNewFromInternalPtr(u).SetIconName(newValue)
				}, widgetPtr)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f Button) Child(widget any) Button {
	return func() *gtk.Button {
		button := f()
		button.SetChild(ResolveWidget(widget))
		return button
	}
}

func (f Button) ConnectClicked(cb func(gtk.Button)) Button {
	return func() *gtk.Button {
		button := f()
		callback.HandleCallback(button.Widget, "clicked", cb)
		return button
	}
}

func (f Button) IconName(iconName string) Button {
	return func() *gtk.Button {
		button := f()
		button.SetIconName(iconName)
		return button
	}
}

func (f Button) Label(label string) Button {
	return func() *gtk.Button {
		button := f()
		button.SetLabel(label)
		return button
	}
}
