package schwifty

import (
	"fmt"

	"codeberg.org/dergs/tidalwave/pkg/schwifty/css"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

type Popover func() *gtk.Popover

func (f Popover) AddController(controller *gtk.EventController) Popover {
 return func() *gtk.Popover {
  widget := f()
  widget.AddController(controller)
  return widget
 }
}

func (f Popover) Background(color string) Popover {
 return func() *gtk.Popover {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { background-color: %s; }", widget.GetCssName(), color))
  return widget
 }
}

func (f Popover) CornerRadius(radius int) Popover {
 return func() *gtk.Popover {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { border-radius: %dpx; }", widget.GetCssName(), radius))
  return widget
 }
}

func (f Popover) CSS(css string) Popover {
	return func() *gtk.Popover {
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

func (f Popover) Focusable(focusable bool) Popover {
	return func() *gtk.Popover {
		widget := f()
		widget.SetFocusable(focusable)
		return widget
	}
}

func (f Popover) FocusOnClick(focusOnClick bool) Popover {
	return func() *gtk.Popover {
		widget := f()
		widget.SetFocusOnClick(focusOnClick)
		return widget
	}
}

func (f Popover) HAlign(align gtk.Align) Popover {
	return func() *gtk.Popover {
		widget := f()
		widget.SetHalign(align)
		return widget
	}
}

func (f Popover) HExpand(expand bool) Popover {
	return func() *gtk.Popover {
		widget := f()
		widget.SetHexpand(expand)
		return widget
	}
}

func (f Popover) HMargin(horizontal int) Popover {
	return func() *gtk.Popover {
		widget := f()
		widget.SetMarginEnd(horizontal)
		widget.SetMarginStart(horizontal)
		return widget
	}
}

func (f Popover) HPadding(padding int) Popover {
 return func() *gtk.Popover {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-left: %dpx; padding-right: %dpx; }", widget.GetCssName(), padding, padding))
  return widget
 }
}

func (f Popover) Margin(margin int) Popover {
	return func() *gtk.Popover {
		widget := f()
		widget.SetMarginBottom(margin)
		widget.SetMarginEnd(margin)
		widget.SetMarginStart(margin)
		widget.SetMarginTop(margin)
		return widget
	}
}

func (f Popover) MarginBottom(bottom int) Popover {
	return func() *gtk.Popover {
		widget := f()
		widget.SetMarginBottom(bottom)
		return widget
	}
}

func (f Popover) MarginEnd(end int) Popover {
	return func() *gtk.Popover {
		widget := f()
		widget.SetMarginEnd(end)
		return widget
	}
}

func (f Popover) MarginStart(start int) Popover {
	return func() *gtk.Popover {
		widget := f()
		widget.SetMarginStart(start)
		return widget
	}
}

func (f Popover) MarginTop(top int) Popover {
	return func() *gtk.Popover {
		widget := f()
		widget.SetMarginTop(top)
		return widget
	}
}

func (f Popover) MinHeight(minHeight int) Popover {
 return func() *gtk.Popover {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { min-height: %dpx; }", widget.GetCssName(), minHeight))
  return widget
 }
}

func (f Popover) MinWidth(minWidth int) Popover {
 return func() *gtk.Popover {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { min-width: %dpx; }", widget.GetCssName(), minWidth))
  return widget
 }
}

func (f Popover) Opacity(opacity float64) Popover {
 return func() *gtk.Popover {
  widget := f()
  widget.SetOpacity(opacity)
  return widget
 }
}

func (f Popover) Overflow(overflow gtk.Overflow) Popover {
 return func() *gtk.Popover {
  widget := f()
  widget.SetOverflow(overflow)
  return widget
 }
}

func (f Popover) Padding(padding int) Popover {
 return func() *gtk.Popover {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { padding: %dpx; }", widget.GetCssName(), padding))
  return widget
 }
}

func (f Popover) PaddingBottom(padding int) Popover {
 return func() *gtk.Popover {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-bottom: %dpx; }", widget.GetCssName(), padding))
  return widget
 }
}

func (f Popover) PaddingEnd(padding int) Popover {
 return func() *gtk.Popover {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-right: %dpx; }", widget.GetCssName(), padding))
  return widget
 }
}

func (f Popover) PaddingStart(padding int) Popover {
 return func() *gtk.Popover {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-left: %dpx; }", widget.GetCssName(), padding))
  return widget
 }
}

func (f Popover) PaddingTop(padding int) Popover {
 return func() *gtk.Popover {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-top: %dpx; }", widget.GetCssName(), padding))
  return widget
 }
}

func (f Popover) ToGTK() *gtk.Widget {
	val := f()
	return &val.Widget
}

func (f Popover) VAlign(align gtk.Align) Popover {
 return func() *gtk.Popover {
  widget := f()
  widget.SetValign(align)
  return widget
 }
}

func (f Popover) VExpand(expand bool) Popover {
 return func() *gtk.Popover {
  widget := f()
  widget.SetVexpand(expand)
  return widget
 }
}

func (f Popover) Visible(visible bool) Popover {
 return func() *gtk.Popover {
  widget := f()
  widget.SetVisible(visible)
  return widget
 }
}

func (f Popover) VMargin(vertical int) Popover {
 return func() *gtk.Popover {
  widget := f()
  widget.SetMarginTop(vertical)
  widget.SetMarginBottom(vertical)
  return widget
 }
}

func (f Popover) VPadding(padding int) Popover {
 return func() *gtk.Popover {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-bottom: %dpx; padding-top: %dpx; }", widget.GetCssName(), padding, padding))
  return widget
 }
}

