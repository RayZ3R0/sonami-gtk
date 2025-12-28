package schwifty

import (
	"codeberg.org/dergs/tidalwave/internal/g"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/css"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	"fmt"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)


type Clamp func() *adw.Clamp

func (f Clamp) AddController(controller *gtk.EventController) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.AddController(controller)
		return widget
	}
}

func (f Clamp) ConnectConstruct(cb func(*adw.Clamp)) Clamp {
	return func() *adw.Clamp {
		widget := f()
		cb(widget)
		return widget
	}
}

func (f Clamp) ConnectDestroy(cb func(gtk.Widget)) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.ConnectDestroy(&cb)
		return widget
	}
}

func (f Clamp) Focusable(focusable bool) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.SetFocusable(focusable)
		return widget
	}
}

func (f Clamp) FocusOnClick(focusOnClick bool) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.SetFocusOnClick(focusOnClick)
		return widget
	}
}

func (f Clamp) HAlign(align gtk.Align) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.SetHalign(align)
		return widget
	}
}

func (f Clamp) HExpand(expand bool) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.SetHexpand(expand)
		return widget
	}
}

func (f Clamp) HMargin(horizontal int) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.SetMarginEnd(horizontal)
		widget.SetMarginStart(horizontal)
		return widget
	}
}

func (f Clamp) Margin(margin int) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.SetMarginBottom(margin)
		widget.SetMarginEnd(margin)
		widget.SetMarginStart(margin)
		widget.SetMarginTop(margin)
		return widget
	}
}

func (f Clamp) MarginBottom(bottom int) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.SetMarginBottom(bottom)
		return widget
	}
}

func (f Clamp) MarginEnd(end int) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.SetMarginEnd(end)
		return widget
	}
}

func (f Clamp) MarginStart(start int) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.SetMarginStart(start)
		return widget
	}
}

func (f Clamp) MarginTop(top int) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.SetMarginTop(top)
		return widget
	}
}

func (f Clamp) Opacity(opacity float64) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.SetOpacity(opacity)
		return widget
	}
}

func (f Clamp) Overflow(overflow gtk.Overflow) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.SetOverflow(overflow)
		return widget
	}
}

func (f Clamp) ToGTK() *gtk.Widget {
	val := f()
	return &val.Widget
}

func (f Clamp) VAlign(align gtk.Align) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.SetValign(align)
		return widget
	}
}

func (f Clamp) VExpand(expand bool) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.SetVexpand(expand)
		return widget
	}
}

func (f Clamp) Visible(visible bool) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.SetVisible(visible)
		return widget
	}
}

func (f Clamp) VMargin(vertical int) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.SetMarginTop(vertical)
		widget.SetMarginBottom(vertical)
		return widget
	}
}



func (f Clamp) Background(color string) Clamp {
	return func() *adw.Clamp {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { background-color: %s; }", widget.GetCssName(), color))
		return widget
	}
}

func (f Clamp) CornerRadius(radius int) Clamp {
	return func() *adw.Clamp {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { border-radius: %dpx; }", widget.GetCssName(), radius))
		return widget
	}
}

func (f Clamp) CSS(css string) Clamp {
	return func() *adw.Clamp {
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

func (f Clamp) HPadding(padding int) Clamp {
	return func() *adw.Clamp {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-left: %dpx; padding-right: %dpx; }", widget.GetCssName(), padding, padding))
		return widget
	}
}

func (f Clamp) MinHeight(minHeight int) Clamp {
	return func() *adw.Clamp {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { min-height: %dpx; }", widget.GetCssName(), minHeight))
		return widget
	}
}

func (f Clamp) MinWidth(minWidth int) Clamp {
	return func() *adw.Clamp {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { min-width: %dpx; }", widget.GetCssName(), minWidth))
		return widget
	}
}

func (f Clamp) Padding(padding int) Clamp {
	return func() *adw.Clamp {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f Clamp) PaddingBottom(padding int) Clamp {
	return func() *adw.Clamp {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-bottom: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f Clamp) PaddingEnd(padding int) Clamp {
	return func() *adw.Clamp {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-right: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f Clamp) PaddingStart(padding int) Clamp {
	return func() *adw.Clamp {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-left: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f Clamp) PaddingTop(padding int) Clamp {
	return func() *adw.Clamp {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-top: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f Clamp) VPadding(padding int) Clamp {
	return func() *adw.Clamp {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-bottom: %dpx; padding-top: %dpx; }", widget.GetCssName(), padding, padding))
		return widget
	}
}



func (f Clamp) BindVisible(state *state.State[bool]) Clamp {
	return func() *adw.Clamp {
		widget := f()

		var callbackId string
		widget.ConnectRealize(g.Ptr(func(a gtk.Widget) {
			callbackId = state.AddCallback(func(newValue bool) {
				a.SetVisible(newValue)
			})
		}))
		widget.ConnectUnrealize(g.Ptr(func(gtk.Widget) {
			state.RemoveCallback(callbackId)
		}))

		return widget
	}
}

func (f Clamp) BindSensitive(state *state.State[bool]) Clamp {
	return func() *adw.Clamp {
		widget := f()

		var callbackId string
		widget.ConnectRealize(g.Ptr(func(a gtk.Widget) {
			callbackId = state.AddCallback(func(newValue bool) {
				a.SetSensitive(newValue)
			})
		}))
		widget.ConnectUnrealize(g.Ptr(func(gtk.Widget) {
			state.RemoveCallback(callbackId)
		}))

		return widget
	}
}
