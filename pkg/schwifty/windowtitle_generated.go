package schwifty

import (
	"codeberg.org/dergs/tidalwave/internal/g"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/css"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	"fmt"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)


type WindowTitle func() *adw.WindowTitle

func (f WindowTitle) AddController(controller *gtk.EventController) WindowTitle {
	return func() *adw.WindowTitle {
		widget := f()
		widget.AddController(controller)
		return widget
	}
}

func (f WindowTitle) ConnectConstruct(cb func(*adw.WindowTitle)) WindowTitle {
	return func() *adw.WindowTitle {
		widget := f()
		cb(widget)
		return widget
	}
}

func (f WindowTitle) ConnectDestroy(cb func(gtk.Widget)) WindowTitle {
	return func() *adw.WindowTitle {
		widget := f()
		widget.ConnectDestroy(&cb)
		return widget
	}
}

func (f WindowTitle) Focusable(focusable bool) WindowTitle {
	return func() *adw.WindowTitle {
		widget := f()
		widget.SetFocusable(focusable)
		return widget
	}
}

func (f WindowTitle) FocusOnClick(focusOnClick bool) WindowTitle {
	return func() *adw.WindowTitle {
		widget := f()
		widget.SetFocusOnClick(focusOnClick)
		return widget
	}
}

func (f WindowTitle) HAlign(align gtk.Align) WindowTitle {
	return func() *adw.WindowTitle {
		widget := f()
		widget.SetHalign(align)
		return widget
	}
}

func (f WindowTitle) HExpand(expand bool) WindowTitle {
	return func() *adw.WindowTitle {
		widget := f()
		widget.SetHexpand(expand)
		return widget
	}
}

func (f WindowTitle) HMargin(horizontal int) WindowTitle {
	return func() *adw.WindowTitle {
		widget := f()
		widget.SetMarginEnd(horizontal)
		widget.SetMarginStart(horizontal)
		return widget
	}
}

func (f WindowTitle) Margin(margin int) WindowTitle {
	return func() *adw.WindowTitle {
		widget := f()
		widget.SetMarginBottom(margin)
		widget.SetMarginEnd(margin)
		widget.SetMarginStart(margin)
		widget.SetMarginTop(margin)
		return widget
	}
}

func (f WindowTitle) MarginBottom(bottom int) WindowTitle {
	return func() *adw.WindowTitle {
		widget := f()
		widget.SetMarginBottom(bottom)
		return widget
	}
}

func (f WindowTitle) MarginEnd(end int) WindowTitle {
	return func() *adw.WindowTitle {
		widget := f()
		widget.SetMarginEnd(end)
		return widget
	}
}

func (f WindowTitle) MarginStart(start int) WindowTitle {
	return func() *adw.WindowTitle {
		widget := f()
		widget.SetMarginStart(start)
		return widget
	}
}

func (f WindowTitle) MarginTop(top int) WindowTitle {
	return func() *adw.WindowTitle {
		widget := f()
		widget.SetMarginTop(top)
		return widget
	}
}

func (f WindowTitle) Opacity(opacity float64) WindowTitle {
	return func() *adw.WindowTitle {
		widget := f()
		widget.SetOpacity(opacity)
		return widget
	}
}

func (f WindowTitle) Overflow(overflow gtk.Overflow) WindowTitle {
	return func() *adw.WindowTitle {
		widget := f()
		widget.SetOverflow(overflow)
		return widget
	}
}

func (f WindowTitle) ToGTK() *gtk.Widget {
	val := f()
	return &val.Widget
}

func (f WindowTitle) VAlign(align gtk.Align) WindowTitle {
	return func() *adw.WindowTitle {
		widget := f()
		widget.SetValign(align)
		return widget
	}
}

func (f WindowTitle) VExpand(expand bool) WindowTitle {
	return func() *adw.WindowTitle {
		widget := f()
		widget.SetVexpand(expand)
		return widget
	}
}

func (f WindowTitle) Visible(visible bool) WindowTitle {
	return func() *adw.WindowTitle {
		widget := f()
		widget.SetVisible(visible)
		return widget
	}
}

func (f WindowTitle) VMargin(vertical int) WindowTitle {
	return func() *adw.WindowTitle {
		widget := f()
		widget.SetMarginTop(vertical)
		widget.SetMarginBottom(vertical)
		return widget
	}
}



func (f WindowTitle) Background(color string) WindowTitle {
	return func() *adw.WindowTitle {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { background-color: %s; }", widget.GetCssName(), color))
		return widget
	}
}

func (f WindowTitle) CornerRadius(radius int) WindowTitle {
	return func() *adw.WindowTitle {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { border-radius: %dpx; }", widget.GetCssName(), radius))
		return widget
	}
}

func (f WindowTitle) CSS(css string) WindowTitle {
	return func() *adw.WindowTitle {
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

func (f WindowTitle) HPadding(padding int) WindowTitle {
	return func() *adw.WindowTitle {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-left: %dpx; padding-right: %dpx; }", widget.GetCssName(), padding, padding))
		return widget
	}
}

func (f WindowTitle) MinHeight(minHeight int) WindowTitle {
	return func() *adw.WindowTitle {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { min-height: %dpx; }", widget.GetCssName(), minHeight))
		return widget
	}
}

func (f WindowTitle) MinWidth(minWidth int) WindowTitle {
	return func() *adw.WindowTitle {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { min-width: %dpx; }", widget.GetCssName(), minWidth))
		return widget
	}
}

func (f WindowTitle) Padding(padding int) WindowTitle {
	return func() *adw.WindowTitle {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f WindowTitle) PaddingBottom(padding int) WindowTitle {
	return func() *adw.WindowTitle {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-bottom: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f WindowTitle) PaddingEnd(padding int) WindowTitle {
	return func() *adw.WindowTitle {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-right: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f WindowTitle) PaddingStart(padding int) WindowTitle {
	return func() *adw.WindowTitle {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-left: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f WindowTitle) PaddingTop(padding int) WindowTitle {
	return func() *adw.WindowTitle {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-top: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f WindowTitle) VPadding(padding int) WindowTitle {
	return func() *adw.WindowTitle {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-bottom: %dpx; padding-top: %dpx; }", widget.GetCssName(), padding, padding))
		return widget
	}
}



func (f WindowTitle) BindVisible(state *state.State[bool]) WindowTitle {
	return func() *adw.WindowTitle {
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

func (f WindowTitle) BindSensitive(state *state.State[bool]) WindowTitle {
	return func() *adw.WindowTitle {
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
