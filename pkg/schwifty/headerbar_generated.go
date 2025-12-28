package schwifty

import (
	"fmt"

	"codeberg.org/dergs/tidalwave/pkg/schwifty/css"
	"github.com/jwijenbergh/puregotk/v4/gtk"
	"github.com/jwijenbergh/puregotk/v4/adw"
)

type HeaderBar func() *adw.HeaderBar

func (f HeaderBar) AddController(controller *gtk.EventController) HeaderBar {
 return func() *adw.HeaderBar {
  widget := f()
  widget.AddController(controller)
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

