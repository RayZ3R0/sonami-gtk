package schwifty

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty/callback"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/css"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	"fmt"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)


type Spinner func() *gtk.Spinner

func (f Spinner) AddController(controller *gtk.EventController) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.AddController(controller)
		return widget
	}
}

func (f Spinner) ConnectConstruct(cb func(*gtk.Spinner)) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		cb(widget)
		return widget
	}
}

func (f Spinner) ConnectDestroy(cb func(gtk.Widget)) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		callback.HandleCallback(widget.Widget, "destroy", cb)
		return widget
	}
}

func (f Spinner) ConnectRealize(cb func(gtk.Widget)) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		callback.HandleCallback(widget.Widget, "realize", cb)
		return widget
	}
}

func (f Spinner) ConnectUnrealize(cb func(gtk.Widget)) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		callback.HandleCallback(widget.Widget, "unrealize", cb)
		return widget
	}
}

func (f Spinner) Focusable(focusable bool) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.SetFocusable(focusable)
		return widget
	}
}

func (f Spinner) FocusOnClick(focusOnClick bool) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.SetFocusOnClick(focusOnClick)
		return widget
	}
}

func (f Spinner) HAlign(align gtk.Align) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.SetHalign(align)
		return widget
	}
}

func (f Spinner) HExpand(expand bool) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.SetHexpand(expand)
		return widget
	}
}

func (f Spinner) HMargin(horizontal int) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.SetMarginEnd(horizontal)
		widget.SetMarginStart(horizontal)
		return widget
	}
}

func (f Spinner) Margin(margin int) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.SetMarginBottom(margin)
		widget.SetMarginEnd(margin)
		widget.SetMarginStart(margin)
		widget.SetMarginTop(margin)
		return widget
	}
}

func (f Spinner) MarginBottom(bottom int) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.SetMarginBottom(bottom)
		return widget
	}
}

func (f Spinner) MarginEnd(end int) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.SetMarginEnd(end)
		return widget
	}
}

func (f Spinner) MarginStart(start int) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.SetMarginStart(start)
		return widget
	}
}

func (f Spinner) MarginTop(top int) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.SetMarginTop(top)
		return widget
	}
}

func (f Spinner) Opacity(opacity float64) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.SetOpacity(opacity)
		return widget
	}
}

func (f Spinner) Overflow(overflow gtk.Overflow) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.SetOverflow(overflow)
		return widget
	}
}

func (f Spinner) ToGTK() *gtk.Widget {
	val := f()
	return &val.Widget
}

func (f Spinner) VAlign(align gtk.Align) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.SetValign(align)
		return widget
	}
}

func (f Spinner) VExpand(expand bool) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.SetVexpand(expand)
		return widget
	}
}

func (f Spinner) Visible(visible bool) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.SetVisible(visible)
		return widget
	}
}

func (f Spinner) VMargin(vertical int) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.SetMarginTop(vertical)
		widget.SetMarginBottom(vertical)
		return widget
	}
}



func (f Spinner) Background(color string) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { background-color: %s; }", widget.GetCssName(), color))
		return widget
	}
}

func (f Spinner) CornerRadius(radius int) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { border-radius: %dpx; }", widget.GetCssName(), radius))
		return widget
	}
}

func (f Spinner) CSS(css string) Spinner {
	return func() *gtk.Spinner {
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

func (f Spinner) HPadding(padding int) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-left: %dpx; padding-right: %dpx; }", widget.GetCssName(), padding, padding))
		return widget
	}
}

func (f Spinner) MinHeight(minHeight int) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { min-height: %dpx; }", widget.GetCssName(), minHeight))
		return widget
	}
}

func (f Spinner) MinWidth(minWidth int) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { min-width: %dpx; }", widget.GetCssName(), minWidth))
		return widget
	}
}

func (f Spinner) Padding(padding int) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f Spinner) PaddingBottom(padding int) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-bottom: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f Spinner) PaddingEnd(padding int) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-right: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f Spinner) PaddingStart(padding int) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-left: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f Spinner) PaddingTop(padding int) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-top: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f Spinner) VPadding(padding int) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-bottom: %dpx; padding-top: %dpx; }", widget.GetCssName(), padding, padding))
		return widget
	}
}



func (f Spinner) BindVisible(state *state.State[bool]) Spinner {
	return func() *gtk.Spinner {
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

func (f Spinner) BindHMargin(state *state.State[int]) Spinner {
	return func() *gtk.Spinner {
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

func (f Spinner) BindMargin(state *state.State[int]) Spinner {
	return func() *gtk.Spinner {
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

func (f Spinner) BindMarginBottom(state *state.State[int]) Spinner {
	return func() *gtk.Spinner {
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

func (f Spinner) BindMarginEnd(state *state.State[int]) Spinner {
	return func() *gtk.Spinner {
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

func (f Spinner) BindMarginStart(state *state.State[int]) Spinner {
	return func() *gtk.Spinner {
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

func (f Spinner) BindMarginTop(state *state.State[int]) Spinner {
	return func() *gtk.Spinner {
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

func (f Spinner) BindSensitive(state *state.State[bool]) Spinner {
	return func() *gtk.Spinner {
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
