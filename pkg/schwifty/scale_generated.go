package schwifty

import (
	"codeberg.org/dergs/tidalwave/internal/g"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/css"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	"fmt"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)


type Scale func() *gtk.Scale

func (f Scale) AddController(controller *gtk.EventController) Scale {
	return func() *gtk.Scale {
		widget := f()
		widget.AddController(controller)
		return widget
	}
}

func (f Scale) ConnectConstruct(cb func(*gtk.Scale)) Scale {
	return func() *gtk.Scale {
		widget := f()
		cb(widget)
		return widget
	}
}

func (f Scale) ConnectDestroy(cb func(gtk.Widget)) Scale {
	return func() *gtk.Scale {
		widget := f()
		widget.ConnectDestroy(&cb)
		return widget
	}
}

func (f Scale) Focusable(focusable bool) Scale {
	return func() *gtk.Scale {
		widget := f()
		widget.SetFocusable(focusable)
		return widget
	}
}

func (f Scale) FocusOnClick(focusOnClick bool) Scale {
	return func() *gtk.Scale {
		widget := f()
		widget.SetFocusOnClick(focusOnClick)
		return widget
	}
}

func (f Scale) HAlign(align gtk.Align) Scale {
	return func() *gtk.Scale {
		widget := f()
		widget.SetHalign(align)
		return widget
	}
}

func (f Scale) HExpand(expand bool) Scale {
	return func() *gtk.Scale {
		widget := f()
		widget.SetHexpand(expand)
		return widget
	}
}

func (f Scale) HMargin(horizontal int) Scale {
	return func() *gtk.Scale {
		widget := f()
		widget.SetMarginEnd(horizontal)
		widget.SetMarginStart(horizontal)
		return widget
	}
}

func (f Scale) Margin(margin int) Scale {
	return func() *gtk.Scale {
		widget := f()
		widget.SetMarginBottom(margin)
		widget.SetMarginEnd(margin)
		widget.SetMarginStart(margin)
		widget.SetMarginTop(margin)
		return widget
	}
}

func (f Scale) MarginBottom(bottom int) Scale {
	return func() *gtk.Scale {
		widget := f()
		widget.SetMarginBottom(bottom)
		return widget
	}
}

func (f Scale) MarginEnd(end int) Scale {
	return func() *gtk.Scale {
		widget := f()
		widget.SetMarginEnd(end)
		return widget
	}
}

func (f Scale) MarginStart(start int) Scale {
	return func() *gtk.Scale {
		widget := f()
		widget.SetMarginStart(start)
		return widget
	}
}

func (f Scale) MarginTop(top int) Scale {
	return func() *gtk.Scale {
		widget := f()
		widget.SetMarginTop(top)
		return widget
	}
}

func (f Scale) Opacity(opacity float64) Scale {
	return func() *gtk.Scale {
		widget := f()
		widget.SetOpacity(opacity)
		return widget
	}
}

func (f Scale) Overflow(overflow gtk.Overflow) Scale {
	return func() *gtk.Scale {
		widget := f()
		widget.SetOverflow(overflow)
		return widget
	}
}

func (f Scale) ToGTK() *gtk.Widget {
	val := f()
	return &val.Widget
}

func (f Scale) VAlign(align gtk.Align) Scale {
	return func() *gtk.Scale {
		widget := f()
		widget.SetValign(align)
		return widget
	}
}

func (f Scale) VExpand(expand bool) Scale {
	return func() *gtk.Scale {
		widget := f()
		widget.SetVexpand(expand)
		return widget
	}
}

func (f Scale) Visible(visible bool) Scale {
	return func() *gtk.Scale {
		widget := f()
		widget.SetVisible(visible)
		return widget
	}
}

func (f Scale) VMargin(vertical int) Scale {
	return func() *gtk.Scale {
		widget := f()
		widget.SetMarginTop(vertical)
		widget.SetMarginBottom(vertical)
		return widget
	}
}



func (f Scale) Background(color string) Scale {
	return func() *gtk.Scale {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { background-color: %s; }", widget.GetCssName(), color))
		return widget
	}
}

func (f Scale) CornerRadius(radius int) Scale {
	return func() *gtk.Scale {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { border-radius: %dpx; }", widget.GetCssName(), radius))
		return widget
	}
}

func (f Scale) CSS(css string) Scale {
	return func() *gtk.Scale {
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

func (f Scale) HPadding(padding int) Scale {
	return func() *gtk.Scale {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-left: %dpx; padding-right: %dpx; }", widget.GetCssName(), padding, padding))
		return widget
	}
}

func (f Scale) MinHeight(minHeight int) Scale {
	return func() *gtk.Scale {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { min-height: %dpx; }", widget.GetCssName(), minHeight))
		return widget
	}
}

func (f Scale) MinWidth(minWidth int) Scale {
	return func() *gtk.Scale {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { min-width: %dpx; }", widget.GetCssName(), minWidth))
		return widget
	}
}

func (f Scale) Padding(padding int) Scale {
	return func() *gtk.Scale {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f Scale) PaddingBottom(padding int) Scale {
	return func() *gtk.Scale {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-bottom: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f Scale) PaddingEnd(padding int) Scale {
	return func() *gtk.Scale {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-right: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f Scale) PaddingStart(padding int) Scale {
	return func() *gtk.Scale {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-left: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f Scale) PaddingTop(padding int) Scale {
	return func() *gtk.Scale {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-top: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f Scale) VPadding(padding int) Scale {
	return func() *gtk.Scale {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-bottom: %dpx; padding-top: %dpx; }", widget.GetCssName(), padding, padding))
		return widget
	}
}



func (f Scale) BindVisible(state *state.State[bool]) Scale {
	return func() *gtk.Scale {
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

func (f Scale) BindSensitive(state *state.State[bool]) Scale {
	return func() *gtk.Scale {
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
