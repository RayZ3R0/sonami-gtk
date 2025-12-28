package schwifty

import (
	"fmt"

	"codeberg.org/dergs/tidalwave/pkg/schwifty/css"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

type Button func() *gtk.Button

func (f Button) AddController(controller *gtk.EventController) Button {
 return func() *gtk.Button {
  widget := f()
  widget.AddController(controller)
  return widget
 }
}

func (f Button) Background(color string) Button {
 return func() *gtk.Button {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { background-color: %s; }", widget.GetCssName(), color))
  return widget
 }
}

func (f Button) CornerRadius(radius int) Button {
 return func() *gtk.Button {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { border-radius: %dpx; }", widget.GetCssName(), radius))
  return widget
 }
}

func (f Button) CSS(css string) Button {
	return func() *gtk.Button {
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

func (f Button) Focusable(focusable bool) Button {
	return func() *gtk.Button {
		widget := f()
		widget.SetFocusable(focusable)
		return widget
	}
}

func (f Button) FocusOnClick(focusOnClick bool) Button {
	return func() *gtk.Button {
		widget := f()
		widget.SetFocusOnClick(focusOnClick)
		return widget
	}
}

func (f Button) HAlign(align gtk.Align) Button {
	return func() *gtk.Button {
		widget := f()
		widget.SetHalign(align)
		return widget
	}
}

func (f Button) HExpand(expand bool) Button {
	return func() *gtk.Button {
		widget := f()
		widget.SetHexpand(expand)
		return widget
	}
}

func (f Button) HMargin(horizontal int) Button {
	return func() *gtk.Button {
		widget := f()
		widget.SetMarginEnd(horizontal)
		widget.SetMarginStart(horizontal)
		return widget
	}
}

func (f Button) HPadding(padding int) Button {
 return func() *gtk.Button {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-left: %dpx; padding-right: %dpx; }", widget.GetCssName(), padding, padding))
  return widget
 }
}

func (f Button) Margin(margin int) Button {
	return func() *gtk.Button {
		widget := f()
		widget.SetMarginBottom(margin)
		widget.SetMarginEnd(margin)
		widget.SetMarginStart(margin)
		widget.SetMarginTop(margin)
		return widget
	}
}

func (f Button) MarginBottom(bottom int) Button {
	return func() *gtk.Button {
		widget := f()
		widget.SetMarginBottom(bottom)
		return widget
	}
}

func (f Button) MarginEnd(end int) Button {
	return func() *gtk.Button {
		widget := f()
		widget.SetMarginEnd(end)
		return widget
	}
}

func (f Button) MarginStart(start int) Button {
	return func() *gtk.Button {
		widget := f()
		widget.SetMarginStart(start)
		return widget
	}
}

func (f Button) MarginTop(top int) Button {
	return func() *gtk.Button {
		widget := f()
		widget.SetMarginTop(top)
		return widget
	}
}

func (f Button) MinHeight(minHeight int) Button {
 return func() *gtk.Button {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { min-height: %dpx; }", widget.GetCssName(), minHeight))
  return widget
 }
}

func (f Button) MinWidth(minWidth int) Button {
 return func() *gtk.Button {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { min-width: %dpx; }", widget.GetCssName(), minWidth))
  return widget
 }
}

func (f Button) Opacity(opacity float64) Button {
 return func() *gtk.Button {
  widget := f()
  widget.SetOpacity(opacity)
  return widget
 }
}

func (f Button) Overflow(overflow gtk.Overflow) Button {
 return func() *gtk.Button {
  widget := f()
  widget.SetOverflow(overflow)
  return widget
 }
}

func (f Button) Padding(padding int) Button {
 return func() *gtk.Button {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { padding: %dpx; }", widget.GetCssName(), padding))
  return widget
 }
}

func (f Button) PaddingBottom(padding int) Button {
 return func() *gtk.Button {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-bottom: %dpx; }", widget.GetCssName(), padding))
  return widget
 }
}

func (f Button) PaddingEnd(padding int) Button {
 return func() *gtk.Button {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-right: %dpx; }", widget.GetCssName(), padding))
  return widget
 }
}

func (f Button) PaddingStart(padding int) Button {
 return func() *gtk.Button {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-left: %dpx; }", widget.GetCssName(), padding))
  return widget
 }
}

func (f Button) PaddingTop(padding int) Button {
 return func() *gtk.Button {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-top: %dpx; }", widget.GetCssName(), padding))
  return widget
 }
}

func (f Button) ToGTK() *gtk.Widget {
	val := f()
	return &val.Widget
}

func (f Button) VAlign(align gtk.Align) Button {
 return func() *gtk.Button {
  widget := f()
  widget.SetValign(align)
  return widget
 }
}

func (f Button) VExpand(expand bool) Button {
 return func() *gtk.Button {
  widget := f()
  widget.SetVexpand(expand)
  return widget
 }
}

func (f Button) Visible(visible bool) Button {
 return func() *gtk.Button {
  widget := f()
  widget.SetVisible(visible)
  return widget
 }
}

func (f Button) VMargin(vertical int) Button {
 return func() *gtk.Button {
  widget := f()
  widget.SetMarginTop(vertical)
  widget.SetMarginBottom(vertical)
  return widget
 }
}

func (f Button) VPadding(padding int) Button {
 return func() *gtk.Button {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-bottom: %dpx; padding-top: %dpx; }", widget.GetCssName(), padding, padding))
  return widget
 }
}

