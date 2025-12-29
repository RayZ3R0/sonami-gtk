package schwifty

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func (f TEMPLATE_TYPE) BindVisible(state *state.State[bool]) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		var callbackId string
		return f.ConnectRealize(func(w gtk.Widget) {
			callbackId = state.AddCallback(func(newValue bool) {
				w.SetVisible(newValue)
			})
		}).ConnectUnrealize(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f TEMPLATE_TYPE) BindHMargin(state *state.State[int]) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		var callbackId string
		return f.ConnectRealize(func(w gtk.Widget) {
			callbackId = state.AddCallback(func(newValue int) {
				w.SetMarginEnd(newValue)
				w.SetMarginStart(newValue)
			})
		}).ConnectUnrealize(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f TEMPLATE_TYPE) BindMargin(state *state.State[int]) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		var callbackId string
		return f.ConnectRealize(func(widget gtk.Widget) {
			callbackId = state.AddCallback(func(newValue int) {
				widget.SetMarginBottom(newValue)
				widget.SetMarginEnd(newValue)
				widget.SetMarginStart(newValue)
				widget.SetMarginTop(newValue)
			})
		}).ConnectUnrealize(func(gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f TEMPLATE_TYPE) BindMarginBottom(state *state.State[int]) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		var callbackId string
		return f.ConnectRealize(func(w gtk.Widget) {
			callbackId = state.AddCallback(func(newValue int) {
				w.SetMarginBottom(newValue)
			})
		}).ConnectUnrealize(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f TEMPLATE_TYPE) BindMarginEnd(state *state.State[int]) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		var callbackId string
		return f.ConnectRealize(func(w gtk.Widget) {
			callbackId = state.AddCallback(func(newValue int) {
				w.SetMarginEnd(newValue)
			})
		}).ConnectUnrealize(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f TEMPLATE_TYPE) BindMarginStart(state *state.State[int]) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		var callbackId string
		return f.ConnectRealize(func(w gtk.Widget) {
			callbackId = state.AddCallback(func(newValue int) {
				w.SetMarginStart(newValue)
			})
		}).ConnectUnrealize(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f TEMPLATE_TYPE) BindMarginTop(state *state.State[int]) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		var callbackId string
		return f.ConnectRealize(func(w gtk.Widget) {
			callbackId = state.AddCallback(func(newValue int) {
				w.SetMarginTop(newValue)
			})
		}).ConnectUnrealize(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f TEMPLATE_TYPE) BindSensitive(state *state.State[bool]) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		var callbackId string
		return f.ConnectRealize(func(w gtk.Widget) {
			callbackId = state.AddCallback(func(newValue bool) {
				w.SetSensitive(newValue)
			})
		}).ConnectUnrealize(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}
