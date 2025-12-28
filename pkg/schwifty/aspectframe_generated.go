package schwifty

import (
	"codeberg.org/dergs/tidalwave/internal/g"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/css"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	"fmt"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)


type AspectFrame func() *gtk.AspectFrame

func (f AspectFrame) AddController(controller *gtk.EventController) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.AddController(controller)
		return widget
	}
}

func (f AspectFrame) ConnectConstruct(cb func(*gtk.AspectFrame)) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		cb(widget)
		return widget
	}
}

func (f AspectFrame) ConnectDestroy(cb func(gtk.Widget)) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.ConnectDestroy(&cb)
		return widget
	}
}

func (f AspectFrame) Focusable(focusable bool) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.SetFocusable(focusable)
		return widget
	}
}

func (f AspectFrame) FocusOnClick(focusOnClick bool) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.SetFocusOnClick(focusOnClick)
		return widget
	}
}

func (f AspectFrame) HAlign(align gtk.Align) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.SetHalign(align)
		return widget
	}
}

func (f AspectFrame) HExpand(expand bool) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.SetHexpand(expand)
		return widget
	}
}

func (f AspectFrame) HMargin(horizontal int) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.SetMarginEnd(horizontal)
		widget.SetMarginStart(horizontal)
		return widget
	}
}

func (f AspectFrame) Margin(margin int) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.SetMarginBottom(margin)
		widget.SetMarginEnd(margin)
		widget.SetMarginStart(margin)
		widget.SetMarginTop(margin)
		return widget
	}
}

func (f AspectFrame) MarginBottom(bottom int) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.SetMarginBottom(bottom)
		return widget
	}
}

func (f AspectFrame) MarginEnd(end int) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.SetMarginEnd(end)
		return widget
	}
}

func (f AspectFrame) MarginStart(start int) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.SetMarginStart(start)
		return widget
	}
}

func (f AspectFrame) MarginTop(top int) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.SetMarginTop(top)
		return widget
	}
}

func (f AspectFrame) Opacity(opacity float64) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.SetOpacity(opacity)
		return widget
	}
}

func (f AspectFrame) Overflow(overflow gtk.Overflow) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.SetOverflow(overflow)
		return widget
	}
}

func (f AspectFrame) ToGTK() *gtk.Widget {
	val := f()
	return &val.Widget
}

func (f AspectFrame) VAlign(align gtk.Align) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.SetValign(align)
		return widget
	}
}

func (f AspectFrame) VExpand(expand bool) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.SetVexpand(expand)
		return widget
	}
}

func (f AspectFrame) Visible(visible bool) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.SetVisible(visible)
		return widget
	}
}

func (f AspectFrame) VMargin(vertical int) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.SetMarginTop(vertical)
		widget.SetMarginBottom(vertical)
		return widget
	}
}



func (f AspectFrame) Background(color string) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { background-color: %s; }", widget.GetCssName(), color))
		return widget
	}
}

func (f AspectFrame) CornerRadius(radius int) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { border-radius: %dpx; }", widget.GetCssName(), radius))
		return widget
	}
}

func (f AspectFrame) CSS(css string) AspectFrame {
	return func() *gtk.AspectFrame {
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

func (f AspectFrame) HPadding(padding int) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-left: %dpx; padding-right: %dpx; }", widget.GetCssName(), padding, padding))
		return widget
	}
}

func (f AspectFrame) MinHeight(minHeight int) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { min-height: %dpx; }", widget.GetCssName(), minHeight))
		return widget
	}
}

func (f AspectFrame) MinWidth(minWidth int) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { min-width: %dpx; }", widget.GetCssName(), minWidth))
		return widget
	}
}

func (f AspectFrame) Padding(padding int) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f AspectFrame) PaddingBottom(padding int) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-bottom: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f AspectFrame) PaddingEnd(padding int) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-right: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f AspectFrame) PaddingStart(padding int) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-left: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f AspectFrame) PaddingTop(padding int) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-top: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f AspectFrame) VPadding(padding int) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-bottom: %dpx; padding-top: %dpx; }", widget.GetCssName(), padding, padding))
		return widget
	}
}



func (f AspectFrame) BindVisible(state *state.State[bool]) AspectFrame {
	return func() *gtk.AspectFrame {
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

func (f AspectFrame) BindSensitive(state *state.State[bool]) AspectFrame {
	return func() *gtk.AspectFrame {
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
