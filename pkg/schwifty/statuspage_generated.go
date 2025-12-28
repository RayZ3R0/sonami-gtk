package schwifty

import (
	"codeberg.org/dergs/tidalwave/internal/g"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/css"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	"fmt"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)


type StatusPage func() *adw.StatusPage

func (f StatusPage) AddController(controller *gtk.EventController) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.AddController(controller)
		return widget
	}
}

func (f StatusPage) ConnectConstruct(cb func(*adw.StatusPage)) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		cb(widget)
		return widget
	}
}

func (f StatusPage) ConnectDestroy(cb func(gtk.Widget)) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.ConnectDestroy(&cb)
		return widget
	}
}

func (f StatusPage) Focusable(focusable bool) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.SetFocusable(focusable)
		return widget
	}
}

func (f StatusPage) FocusOnClick(focusOnClick bool) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.SetFocusOnClick(focusOnClick)
		return widget
	}
}

func (f StatusPage) HAlign(align gtk.Align) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.SetHalign(align)
		return widget
	}
}

func (f StatusPage) HExpand(expand bool) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.SetHexpand(expand)
		return widget
	}
}

func (f StatusPage) HMargin(horizontal int) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.SetMarginEnd(horizontal)
		widget.SetMarginStart(horizontal)
		return widget
	}
}

func (f StatusPage) Margin(margin int) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.SetMarginBottom(margin)
		widget.SetMarginEnd(margin)
		widget.SetMarginStart(margin)
		widget.SetMarginTop(margin)
		return widget
	}
}

func (f StatusPage) MarginBottom(bottom int) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.SetMarginBottom(bottom)
		return widget
	}
}

func (f StatusPage) MarginEnd(end int) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.SetMarginEnd(end)
		return widget
	}
}

func (f StatusPage) MarginStart(start int) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.SetMarginStart(start)
		return widget
	}
}

func (f StatusPage) MarginTop(top int) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.SetMarginTop(top)
		return widget
	}
}

func (f StatusPage) Opacity(opacity float64) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.SetOpacity(opacity)
		return widget
	}
}

func (f StatusPage) Overflow(overflow gtk.Overflow) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.SetOverflow(overflow)
		return widget
	}
}

func (f StatusPage) ToGTK() *gtk.Widget {
	val := f()
	return &val.Widget
}

func (f StatusPage) VAlign(align gtk.Align) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.SetValign(align)
		return widget
	}
}

func (f StatusPage) VExpand(expand bool) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.SetVexpand(expand)
		return widget
	}
}

func (f StatusPage) Visible(visible bool) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.SetVisible(visible)
		return widget
	}
}

func (f StatusPage) VMargin(vertical int) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.SetMarginTop(vertical)
		widget.SetMarginBottom(vertical)
		return widget
	}
}



func (f StatusPage) Background(color string) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { background-color: %s; }", widget.GetCssName(), color))
		return widget
	}
}

func (f StatusPage) CornerRadius(radius int) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { border-radius: %dpx; }", widget.GetCssName(), radius))
		return widget
	}
}

func (f StatusPage) CSS(css string) StatusPage {
	return func() *adw.StatusPage {
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

func (f StatusPage) HPadding(padding int) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-left: %dpx; padding-right: %dpx; }", widget.GetCssName(), padding, padding))
		return widget
	}
}

func (f StatusPage) MinHeight(minHeight int) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { min-height: %dpx; }", widget.GetCssName(), minHeight))
		return widget
	}
}

func (f StatusPage) MinWidth(minWidth int) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { min-width: %dpx; }", widget.GetCssName(), minWidth))
		return widget
	}
}

func (f StatusPage) Padding(padding int) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f StatusPage) PaddingBottom(padding int) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-bottom: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f StatusPage) PaddingEnd(padding int) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-right: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f StatusPage) PaddingStart(padding int) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-left: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f StatusPage) PaddingTop(padding int) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-top: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f StatusPage) VPadding(padding int) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-bottom: %dpx; padding-top: %dpx; }", widget.GetCssName(), padding, padding))
		return widget
	}
}



func (f StatusPage) BindVisible(state *state.State[bool]) StatusPage {
	return func() *adw.StatusPage {
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

func (f StatusPage) BindSensitive(state *state.State[bool]) StatusPage {
	return func() *adw.StatusPage {
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
