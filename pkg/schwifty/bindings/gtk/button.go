package gtk

import (
	"codeberg.org/puregotk/puregotk/v4/glib"
	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/callback"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/state"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/utils/weak"
)

//go:generate go run github.com/RayZ3R0/sonami-gtk/pkg/schwifty/gen Button *gtk.Button gtk

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
		var ref weak.WidgetRef
		return f.ConnectRealize(func(w gtk.Widget) {
			ref = weak.NewWidgetRef(&w)
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
		}).ConnectUnrealize(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f Button) BindIconName(state *state.State[string]) Button {
	return func() *gtk.Button {
		var callbackId string
		var ref weak.WidgetRef
		return f.ConnectRealize(func(w gtk.Widget) {
			ref = weak.NewWidgetRef(&w)
			callbackId = state.AddCallback(func(newValue string) {
				callback.OnMainThreadOncePure(func() {
					if obj := ref.Get(); obj != nil {
						defer obj.Unref()
						gtk.ButtonNewFromInternalPtr(obj.Ptr).SetIconName(newValue)
					}
				})
			})
		}).ConnectUnrealize(func(w gtk.Widget) {
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

func (f Button) BindTooltipText(state *state.State[string]) Button {
	return func() *gtk.Button {
		var callbackId string
		var ref weak.WidgetRef
		return f.ConnectRealize(func(w gtk.Widget) {
			ref = weak.NewWidgetRef(&w)
			callbackId = state.AddCallback(func(newValue string) {
				callback.OnMainThreadOncePure(func() {
					if obj := ref.Get(); obj != nil {
						defer obj.Unref()
						gtk.ButtonNewFromInternalPtr(obj.Ptr).SetTooltipText(newValue)
					}
				})
			})
		}).ConnectUnrealize(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}
