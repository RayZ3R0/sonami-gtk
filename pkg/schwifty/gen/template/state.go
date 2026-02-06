package schwifty

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty/callback"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	"codeberg.org/dergs/tonearm/pkg/schwifty/tracking"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func (f TEMPLATE_TYPE) BindVisible(state *state.State[bool]) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		var callbackId string
		var ref *tracking.WeakRef
		return f.ConnectConstruct(func(w TEMPLATE_BASE_TYPE) {
			ref = tracking.NewWeakRef(&w.Widget)
			callbackId = state.AddCallback(func(newValue bool) {
				callback.OnMainThreadOncePure(func() {
					if obj := ref.Get(); obj != nil {
						defer obj.Unref()
						gtk.WidgetNewFromInternalPtr(obj.Ptr).SetVisible(newValue)
					}
				})
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f TEMPLATE_TYPE) BindHMargin(state *state.State[int]) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		var callbackId string
		var ref *tracking.WeakRef
		return f.ConnectConstruct(func(w TEMPLATE_BASE_TYPE) {
			ref = tracking.NewWeakRef(&w.Widget)
			callbackId = state.AddCallback(func(newValue int) {
				callback.OnMainThreadOncePure(func() {
					if obj := ref.Get(); obj != nil {
						defer obj.Unref()
						gtk.WidgetNewFromInternalPtr(obj.Ptr).SetMarginEnd(newValue)
						gtk.WidgetNewFromInternalPtr(obj.Ptr).SetMarginStart(newValue)
					}
				})
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f TEMPLATE_TYPE) BindMargin(state *state.State[int]) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		var callbackId string
		var ref *tracking.WeakRef
		return f.ConnectConstruct(func(w TEMPLATE_BASE_TYPE) {
			ref = tracking.NewWeakRef(&w.Widget)
			callbackId = state.AddCallback(func(newValue int) {
				callback.OnMainThreadOncePure(func() {
					if obj := ref.Get(); obj != nil {
						defer obj.Unref()
						gtk.WidgetNewFromInternalPtr(obj.Ptr).SetMarginBottom(newValue)
						gtk.WidgetNewFromInternalPtr(obj.Ptr).SetMarginEnd(newValue)
						gtk.WidgetNewFromInternalPtr(obj.Ptr).SetMarginStart(newValue)
						gtk.WidgetNewFromInternalPtr(obj.Ptr).SetMarginTop(newValue)
					}
				})
			})
		}).ConnectDestroy(func(gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f TEMPLATE_TYPE) BindMarginBottom(state *state.State[int]) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		var callbackId string
		var ref *tracking.WeakRef
		return f.ConnectConstruct(func(w TEMPLATE_BASE_TYPE) {
			ref = tracking.NewWeakRef(&w.Widget)
			callbackId = state.AddCallback(func(newValue int) {
				callback.OnMainThreadOncePure(func() {
					if obj := ref.Get(); obj != nil {
						defer obj.Unref()
						gtk.WidgetNewFromInternalPtr(obj.Ptr).SetMarginBottom(newValue)
					}
				})
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f TEMPLATE_TYPE) BindMarginEnd(state *state.State[int]) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		var callbackId string
		var ref *tracking.WeakRef
		return f.ConnectConstruct(func(w TEMPLATE_BASE_TYPE) {
			ref = tracking.NewWeakRef(&w.Widget)
			callbackId = state.AddCallback(func(newValue int) {
				callback.OnMainThreadOncePure(func() {
					if obj := ref.Get(); obj != nil {
						defer obj.Unref()
						gtk.WidgetNewFromInternalPtr(obj.Ptr).SetMarginEnd(newValue)
					}
				})
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f TEMPLATE_TYPE) BindMarginStart(state *state.State[int]) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		var callbackId string
		var ref *tracking.WeakRef
		return f.ConnectConstruct(func(w TEMPLATE_BASE_TYPE) {
			ref = tracking.NewWeakRef(&w.Widget)
			callbackId = state.AddCallback(func(newValue int) {
				callback.OnMainThreadOncePure(func() {
					if obj := ref.Get(); obj != nil {
						defer obj.Unref()
						gtk.WidgetNewFromInternalPtr(obj.Ptr).SetMarginStart(newValue)
					}
				})
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f TEMPLATE_TYPE) BindMarginTop(state *state.State[int]) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		var callbackId string
		var ref *tracking.WeakRef
		return f.ConnectConstruct(func(w TEMPLATE_BASE_TYPE) {
			ref = tracking.NewWeakRef(&w.Widget)
			callbackId = state.AddCallback(func(newValue int) {
				callback.OnMainThreadOncePure(func() {
					if obj := ref.Get(); obj != nil {
						defer obj.Unref()
						gtk.WidgetNewFromInternalPtr(obj.Ptr).SetMarginTop(newValue)
					}
				})
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f TEMPLATE_TYPE) BindSensitive(state *state.State[bool]) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		var callbackId string
		var ref *tracking.WeakRef
		return f.ConnectConstruct(func(w TEMPLATE_BASE_TYPE) {
			ref = tracking.NewWeakRef(&w.Widget)
			callbackId = state.AddCallback(func(newValue bool) {
				callback.OnMainThreadOncePure(func() {
					if obj := ref.Get(); obj != nil {
						defer obj.Unref()
						gtk.WidgetNewFromInternalPtr(obj.Ptr).SetSensitive(newValue)
					}
				})
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}
