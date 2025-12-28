package schwifty

import (
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
		button.ConnectClicked(&cb)
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
