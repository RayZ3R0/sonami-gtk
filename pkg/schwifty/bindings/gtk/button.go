package gtk

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty/callback"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	"codeberg.org/dergs/tonearm/pkg/schwifty/tracking"
	"github.com/jwijenbergh/puregotk/v4/glib"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

//go:generate go run codeberg.org/dergs/tonearm/pkg/schwifty/gen Button *gtk.Button gtk

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

func (f Button) BindChild(state *state.State[any]) Button {
	return func() *gtk.Button {
		var callbackId string
		var ref *tracking.WeakRef
		return f.ConnectConstruct(func(w *gtk.Button) {
			ref = tracking.NewWeakRef(w)
			callbackId = state.AddCallback(func(newValue any) {
				widget := ResolveWidget(newValue)
				if widget == nil {
					callback.OnMainThreadOncePure(func() {
						if obj := ref.Get(); obj != nil {
							defer obj.Unref()
							gtk.ButtonNewFromInternalPtr(obj.Ptr).SetChild(nil)
						}
					})
				} else {
					widget.Ref()
					callback.OnMainThreadOncePure(func() {
						defer widget.Unref()
						if obj := ref.Get(); obj != nil {
							defer obj.Unref()
							gtk.ButtonNewFromInternalPtr(obj.Ptr).SetChild(widget)
						}
					})
				}
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f Button) BindIconName(state *state.State[string]) Button {
	return func() *gtk.Button {
		var callbackId string
		var ref *tracking.WeakRef
		return f.ConnectConstruct(func(w *gtk.Button) {
			ref = tracking.NewWeakRef(w)
			callbackId = state.AddCallback(func(newValue string) {
				callback.OnMainThreadOncePure(func() {
					if obj := ref.Get(); obj != nil {
						defer obj.Unref()
						gtk.ButtonNewFromInternalPtr(obj.Ptr).SetIconName(newValue)
					}
				})
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
		callback.HandleCallback(button.Object, "clicked", cb)
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

func (f Button) TooltipText(tooltip string) Button {
	return func() *gtk.Button {
		button := f()
		button.SetTooltipText(tooltip)
		return button
	}
}
