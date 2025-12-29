package schwifty

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty/callback"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/css"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	"fmt"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)


type HeaderBar func() *adw.HeaderBar

func (f HeaderBar) AddController(controller *gtk.EventController) HeaderBar {
	return func() *adw.HeaderBar {
		widget := f()
		widget.AddController(controller)
		return widget
	}
}

func (f HeaderBar) ConnectConstruct(cb func(*adw.HeaderBar)) HeaderBar {
	return func() *adw.HeaderBar {
		widget := f()
		cb(widget)
		return widget
	}
}

func (f HeaderBar) ConnectDestroy(cb func(gtk.Widget)) HeaderBar {
	return func() *adw.HeaderBar {
		widget := f()
		callback.HandleCallback(widget.Widget, "destroy", cb)
		return widget
	}
}

func (f HeaderBar) ConnectRealize(cb func(gtk.Widget)) HeaderBar {
	return func() *adw.HeaderBar {
		widget := f()
		callback.HandleCallback(widget.Widget, "realize", cb)
		return widget
	}
}

func (f HeaderBar) ConnectUnrealize(cb func(gtk.Widget)) HeaderBar {
	return func() *adw.HeaderBar {
		widget := f()
		callback.HandleCallback(widget.Widget, "unrealize", cb)
		return widget
	}
}

func (f HeaderBar) Focusable(focusable bool) HeaderBar {
	return func() *adw.HeaderBar {
		widget := f()
		widget.SetFocusable(focusable)
		return widget
	}
}

func (f HeaderBar) FocusOnClick(focusOnClick bool) HeaderBar {
	return func() *adw.HeaderBar {
		widget := f()
		widget.SetFocusOnClick(focusOnClick)
		return widget
	}
}

func (f HeaderBar) HAlign(align gtk.Align) HeaderBar {
	return func() *adw.HeaderBar {
		widget := f()
		widget.SetHalign(align)
		return widget
	}
}

func (f HeaderBar) HExpand(expand bool) HeaderBar {
	return func() *adw.HeaderBar {
		widget := f()
		widget.SetHexpand(expand)
		return widget
	}
}

func (f HeaderBar) HMargin(horizontal int) HeaderBar {
	return func() *adw.HeaderBar {
		widget := f()
		widget.SetMarginEnd(horizontal)
		widget.SetMarginStart(horizontal)
		return widget
	}
}

func (f HeaderBar) Margin(margin int) HeaderBar {
	return func() *adw.HeaderBar {
		widget := f()
		widget.SetMarginBottom(margin)
		widget.SetMarginEnd(margin)
		widget.SetMarginStart(margin)
		widget.SetMarginTop(margin)
		return widget
	}
}

func (f HeaderBar) MarginBottom(bottom int) HeaderBar {
	return func() *adw.HeaderBar {
		widget := f()
		widget.SetMarginBottom(bottom)
		return widget
	}
}

func (f HeaderBar) MarginEnd(end int) HeaderBar {
	return func() *adw.HeaderBar {
		widget := f()
		widget.SetMarginEnd(end)
		return widget
	}
}

func (f HeaderBar) MarginStart(start int) HeaderBar {
	return func() *adw.HeaderBar {
		widget := f()
		widget.SetMarginStart(start)
		return widget
	}
}

func (f HeaderBar) MarginTop(top int) HeaderBar {
	return func() *adw.HeaderBar {
		widget := f()
		widget.SetMarginTop(top)
		return widget
	}
}

func (f HeaderBar) Opacity(opacity float64) HeaderBar {
	return func() *adw.HeaderBar {
		widget := f()
		widget.SetOpacity(opacity)
		return widget
	}
}

func (f HeaderBar) Overflow(overflow gtk.Overflow) HeaderBar {
	return func() *adw.HeaderBar {
		widget := f()
		widget.SetOverflow(overflow)
		return widget
	}
}

func (f HeaderBar) ToGTK() *gtk.Widget {
	val := f()
	return &val.Widget
}

func (f HeaderBar) VAlign(align gtk.Align) HeaderBar {
	return func() *adw.HeaderBar {
		widget := f()
		widget.SetValign(align)
		return widget
	}
}

func (f HeaderBar) VExpand(expand bool) HeaderBar {
	return func() *adw.HeaderBar {
		widget := f()
		widget.SetVexpand(expand)
		return widget
	}
}

func (f HeaderBar) Visible(visible bool) HeaderBar {
	return func() *adw.HeaderBar {
		widget := f()
		widget.SetVisible(visible)
		return widget
	}
}

func (f HeaderBar) VMargin(vertical int) HeaderBar {
	return func() *adw.HeaderBar {
		widget := f()
		widget.SetMarginTop(vertical)
		widget.SetMarginBottom(vertical)
		return widget
	}
}



func (f HeaderBar) Background(color string) HeaderBar {
	return func() *adw.HeaderBar {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { background-color: %s; }", widget.GetCssName(), color))
		return widget
	}
}

func (f HeaderBar) CornerRadius(radius int) HeaderBar {
	return func() *adw.HeaderBar {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { border-radius: %dpx; }", widget.GetCssName(), radius))
		return widget
	}
}

func (f HeaderBar) CSS(css string) HeaderBar {
	return func() *adw.HeaderBar {
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

func (f HeaderBar) HPadding(padding int) HeaderBar {
	return func() *adw.HeaderBar {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-left: %dpx; padding-right: %dpx; }", widget.GetCssName(), padding, padding))
		return widget
	}
}

func (f HeaderBar) MinHeight(minHeight int) HeaderBar {
	return func() *adw.HeaderBar {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { min-height: %dpx; }", widget.GetCssName(), minHeight))
		return widget
	}
}

func (f HeaderBar) MinWidth(minWidth int) HeaderBar {
	return func() *adw.HeaderBar {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { min-width: %dpx; }", widget.GetCssName(), minWidth))
		return widget
	}
}

func (f HeaderBar) Padding(padding int) HeaderBar {
	return func() *adw.HeaderBar {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f HeaderBar) PaddingBottom(padding int) HeaderBar {
	return func() *adw.HeaderBar {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-bottom: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f HeaderBar) PaddingEnd(padding int) HeaderBar {
	return func() *adw.HeaderBar {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-right: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f HeaderBar) PaddingStart(padding int) HeaderBar {
	return func() *adw.HeaderBar {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-left: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f HeaderBar) PaddingTop(padding int) HeaderBar {
	return func() *adw.HeaderBar {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-top: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f HeaderBar) VPadding(padding int) HeaderBar {
	return func() *adw.HeaderBar {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-bottom: %dpx; padding-top: %dpx; }", widget.GetCssName(), padding, padding))
		return widget
	}
}



func (f HeaderBar) BindVisible(state *state.State[bool]) HeaderBar {
	return func() *adw.HeaderBar {
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

func (f HeaderBar) BindHMargin(state *state.State[int]) HeaderBar {
	return func() *adw.HeaderBar {
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

func (f HeaderBar) BindMargin(state *state.State[int]) HeaderBar {
	return func() *adw.HeaderBar {
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

func (f HeaderBar) BindMarginBottom(state *state.State[int]) HeaderBar {
	return func() *adw.HeaderBar {
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

func (f HeaderBar) BindMarginEnd(state *state.State[int]) HeaderBar {
	return func() *adw.HeaderBar {
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

func (f HeaderBar) BindMarginStart(state *state.State[int]) HeaderBar {
	return func() *adw.HeaderBar {
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

func (f HeaderBar) BindMarginTop(state *state.State[int]) HeaderBar {
	return func() *adw.HeaderBar {
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

func (f HeaderBar) BindSensitive(state *state.State[bool]) HeaderBar {
	return func() *adw.HeaderBar {
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
