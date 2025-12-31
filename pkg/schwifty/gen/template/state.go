package schwifty

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func (f TEMPLATE_TYPE) BindVisible(state *state.State[bool]) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		var callbackId string
		return f.ConnectConstruct(func(w TEMPLATE_BASE_TYPE) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue bool) {
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetVisible(newValue)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f TEMPLATE_TYPE) BindHMargin(state *state.State[int]) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		var callbackId string
		return f.ConnectConstruct(func(w TEMPLATE_BASE_TYPE) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue int) {
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetMarginEnd(newValue)
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetMarginStart(newValue)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f TEMPLATE_TYPE) BindMargin(state *state.State[int]) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		var callbackId string
		return f.ConnectConstruct(func(widget TEMPLATE_BASE_TYPE) {
			widgetPtr := widget.GoPointer()
			callbackId = state.AddCallback(func(newValue int) {
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetMarginBottom(newValue)
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetMarginEnd(newValue)
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetMarginStart(newValue)
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetMarginTop(newValue)
			})
		}).ConnectDestroy(func(gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f TEMPLATE_TYPE) BindMarginBottom(state *state.State[int]) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		var callbackId string
		return f.ConnectConstruct(func(w TEMPLATE_BASE_TYPE) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue int) {
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetMarginBottom(newValue)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f TEMPLATE_TYPE) BindMarginEnd(state *state.State[int]) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		var callbackId string
		return f.ConnectConstruct(func(w TEMPLATE_BASE_TYPE) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue int) {
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetMarginEnd(newValue)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f TEMPLATE_TYPE) BindMarginStart(state *state.State[int]) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		var callbackId string
		return f.ConnectConstruct(func(w TEMPLATE_BASE_TYPE) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue int) {
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetMarginStart(newValue)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f TEMPLATE_TYPE) BindMarginTop(state *state.State[int]) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		var callbackId string
		return f.ConnectConstruct(func(w TEMPLATE_BASE_TYPE) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue int) {
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetMarginTop(newValue)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f TEMPLATE_TYPE) BindSensitive(state *state.State[bool]) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		var callbackId string
		return f.ConnectConstruct(func(w TEMPLATE_BASE_TYPE) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue bool) {
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetSensitive(newValue)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}
