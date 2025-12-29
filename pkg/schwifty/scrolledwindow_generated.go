package schwifty

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty/callback"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/css"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	"fmt"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)


type ScrolledWindow func() *gtk.ScrolledWindow

func (f ScrolledWindow) AddController(controller *gtk.EventController) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.AddController(controller)
		return widget
	}
}

func (f ScrolledWindow) ConnectConstruct(cb func(*gtk.ScrolledWindow)) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		cb(widget)
		return widget
	}
}

func (f ScrolledWindow) ConnectDestroy(cb func(gtk.Widget)) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		callback.HandleCallback(widget.Widget, "destroy", cb)
		return widget
	}
}

func (f ScrolledWindow) ConnectRealize(cb func(gtk.Widget)) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		callback.HandleCallback(widget.Widget, "realize", cb)
		return widget
	}
}

func (f ScrolledWindow) ConnectUnrealize(cb func(gtk.Widget)) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		callback.HandleCallback(widget.Widget, "unrealize", cb)
		return widget
	}
}

func (f ScrolledWindow) Focusable(focusable bool) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.SetFocusable(focusable)
		return widget
	}
}

func (f ScrolledWindow) FocusOnClick(focusOnClick bool) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.SetFocusOnClick(focusOnClick)
		return widget
	}
}

func (f ScrolledWindow) HAlign(align gtk.Align) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.SetHalign(align)
		return widget
	}
}

func (f ScrolledWindow) HExpand(expand bool) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.SetHexpand(expand)
		return widget
	}
}

func (f ScrolledWindow) HMargin(horizontal int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.SetMarginEnd(horizontal)
		widget.SetMarginStart(horizontal)
		return widget
	}
}

func (f ScrolledWindow) Margin(margin int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.SetMarginBottom(margin)
		widget.SetMarginEnd(margin)
		widget.SetMarginStart(margin)
		widget.SetMarginTop(margin)
		return widget
	}
}

func (f ScrolledWindow) MarginBottom(bottom int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.SetMarginBottom(bottom)
		return widget
	}
}

func (f ScrolledWindow) MarginEnd(end int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.SetMarginEnd(end)
		return widget
	}
}

func (f ScrolledWindow) MarginStart(start int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.SetMarginStart(start)
		return widget
	}
}

func (f ScrolledWindow) MarginTop(top int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.SetMarginTop(top)
		return widget
	}
}

func (f ScrolledWindow) Opacity(opacity float64) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.SetOpacity(opacity)
		return widget
	}
}

func (f ScrolledWindow) Overflow(overflow gtk.Overflow) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.SetOverflow(overflow)
		return widget
	}
}

func (f ScrolledWindow) ToGTK() *gtk.Widget {
	val := f()
	return &val.Widget
}

func (f ScrolledWindow) VAlign(align gtk.Align) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.SetValign(align)
		return widget
	}
}

func (f ScrolledWindow) VExpand(expand bool) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.SetVexpand(expand)
		return widget
	}
}

func (f ScrolledWindow) Visible(visible bool) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.SetVisible(visible)
		return widget
	}
}

func (f ScrolledWindow) VMargin(vertical int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.SetMarginTop(vertical)
		widget.SetMarginBottom(vertical)
		return widget
	}
}



func (f ScrolledWindow) Background(color string) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { background-color: %s; }", widget.GetCssName(), color))
		return widget
	}
}

func (f ScrolledWindow) CornerRadius(radius int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { border-radius: %dpx; }", widget.GetCssName(), radius))
		return widget
	}
}

func (f ScrolledWindow) CSS(css string) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.Ref()
		defer widget.Unref()

		provider := gtk.NewCssProvider()
		provider.LoadFromString(css)
		widget.GetStyleContext().AddProvider(provider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))
		provider.Unref()

		return widget
	}
}

func (f ScrolledWindow) HPadding(padding int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-left: %dpx; padding-right: %dpx; }", widget.GetCssName(), padding, padding))
		return widget
	}
}

func (f ScrolledWindow) MinHeight(minHeight int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { min-height: %dpx; }", widget.GetCssName(), minHeight))
		return widget
	}
}

func (f ScrolledWindow) MinWidth(minWidth int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { min-width: %dpx; }", widget.GetCssName(), minWidth))
		return widget
	}
}

func (f ScrolledWindow) Padding(padding int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f ScrolledWindow) PaddingBottom(padding int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-bottom: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f ScrolledWindow) PaddingEnd(padding int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-right: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f ScrolledWindow) PaddingStart(padding int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-left: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f ScrolledWindow) PaddingTop(padding int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-top: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f ScrolledWindow) VPadding(padding int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-bottom: %dpx; padding-top: %dpx; }", widget.GetCssName(), padding, padding))
		return widget
	}
}



func (f ScrolledWindow) BindVisible(state *state.State[bool]) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
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

func (f ScrolledWindow) BindHMargin(state *state.State[int]) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
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

func (f ScrolledWindow) BindMargin(state *state.State[int]) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
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

func (f ScrolledWindow) BindMarginBottom(state *state.State[int]) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
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

func (f ScrolledWindow) BindMarginEnd(state *state.State[int]) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
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

func (f ScrolledWindow) BindMarginStart(state *state.State[int]) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
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

func (f ScrolledWindow) BindMarginTop(state *state.State[int]) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
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

func (f ScrolledWindow) BindSensitive(state *state.State[bool]) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
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
