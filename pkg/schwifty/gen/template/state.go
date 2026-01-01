package schwifty

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty/callback"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func (f TEMPLATE_TYPE) BindVisible(state *state.State[bool]) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		var callbackId string
		return f.ConnectConstruct(func(w TEMPLATE_BASE_TYPE) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue bool) {
				callback.OnMainThreadOnce(func(u uintptr) {
					gtk.WidgetNewFromInternalPtr(u).SetVisible(newValue)
				}, widgetPtr)
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
				callback.OnMainThreadOnce(func(u uintptr) {
					gtk.WidgetNewFromInternalPtr(u).SetMarginEnd(newValue)
					gtk.WidgetNewFromInternalPtr(u).SetMarginStart(newValue)
				}, widgetPtr)
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
				callback.OnMainThreadOnce(func(u uintptr) {
					gtk.WidgetNewFromInternalPtr(u).SetMarginBottom(newValue)
					gtk.WidgetNewFromInternalPtr(u).SetMarginEnd(newValue)
					gtk.WidgetNewFromInternalPtr(u).SetMarginStart(newValue)
					gtk.WidgetNewFromInternalPtr(u).SetMarginTop(newValue)
				}, widgetPtr)
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
				callback.OnMainThreadOnce(func(u uintptr) {
					gtk.WidgetNewFromInternalPtr(u).SetMarginBottom(newValue)
				}, widgetPtr)
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
				callback.OnMainThreadOnce(func(u uintptr) {
					gtk.WidgetNewFromInternalPtr(u).SetMarginEnd(newValue)
				}, widgetPtr)
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
				callback.OnMainThreadOnce(func(u uintptr) {
					gtk.WidgetNewFromInternalPtr(u).SetMarginStart(newValue)
				}, widgetPtr)
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
				callback.OnMainThreadOnce(func(u uintptr) {
					gtk.WidgetNewFromInternalPtr(u).SetMarginTop(newValue)
				}, widgetPtr)
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
				callback.OnMainThreadOnce(func(u uintptr) {
					gtk.WidgetNewFromInternalPtr(u).SetSensitive(newValue)
				}, widgetPtr)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}
