package schwifty

import (
	"fmt"

	"codeberg.org/dergs/tidalwave/pkg/schwifty/css"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

type Box func() *gtk.Box

func (f Box) AddController(controller *gtk.EventController) Box {
 return func() *gtk.Box {
  widget := f()
  widget.AddController(controller)
  return widget
 }
}

func (f Box) Background(color string) Box {
 return func() *gtk.Box {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { background-color: %s; }", widget.GetCssName(), color))
  return widget
 }
}

func (f Box) CornerRadius(radius int) Box {
 return func() *gtk.Box {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { border-radius: %dpx; }", widget.GetCssName(), radius))
  return widget
 }
}

func (f Box) CSS(css string) Box {
	return func() *gtk.Box {
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

func (f Box) Focusable(focusable bool) Box {
	return func() *gtk.Box {
		widget := f()
		widget.SetFocusable(focusable)
		return widget
	}
}

func (f Box) FocusOnClick(focusOnClick bool) Box {
	return func() *gtk.Box {
		widget := f()
		widget.SetFocusOnClick(focusOnClick)
		return widget
	}
}

func (f Box) HAlign(align gtk.Align) Box {
	return func() *gtk.Box {
		widget := f()
		widget.SetHalign(align)
		return widget
	}
}

func (f Box) HExpand(expand bool) Box {
	return func() *gtk.Box {
		widget := f()
		widget.SetHexpand(expand)
		return widget
	}
}

func (f Box) HMargin(horizontal int) Box {
	return func() *gtk.Box {
		widget := f()
		widget.SetMarginEnd(horizontal)
		widget.SetMarginStart(horizontal)
		return widget
	}
}

func (f Box) Margin(margin int) Box {
	return func() *gtk.Box {
		widget := f()
		widget.SetMarginBottom(margin)
		widget.SetMarginEnd(margin)
		widget.SetMarginStart(margin)
		widget.SetMarginTop(margin)
		return widget
	}
}

func (f Box) MarginBottom(bottom int) Box {
	return func() *gtk.Box {
		widget := f()
		widget.SetMarginBottom(bottom)
		return widget
	}
}

func (f Box) MarginEnd(end int) Box {
	return func() *gtk.Box {
		widget := f()
		widget.SetMarginEnd(end)
		return widget
	}
}

func (f Box) MarginStart(start int) Box {
	return func() *gtk.Box {
		widget := f()
		widget.SetMarginStart(start)
		return widget
	}
}

func (f Box) MarginTop(top int) Box {
	return func() *gtk.Box {
		widget := f()
		widget.SetMarginTop(top)
		return widget
	}
}

func (f Box) MinHeight(minHeight int) Box {
 return func() *gtk.Box {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { min-height: %dpx; }", widget.GetCssName(), minHeight))
  return widget
 }
}

func (f Box) MinWidth(minWidth int) Box {
 return func() *gtk.Box {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { min-width: %dpx; }", widget.GetCssName(), minWidth))
  return widget
 }
}

func (f Box) Opacity(opacity float64) Box {
 return func() *gtk.Box {
  widget := f()
  widget.SetOpacity(opacity)
  return widget
 }
}

func (f Box) Overflow(overflow gtk.Overflow) Box {
 return func() *gtk.Box {
  widget := f()
  widget.SetOverflow(overflow)
  return widget
 }
}

func (f Box) Padding(padding int) Box {
 return func() *gtk.Box {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { padding: %dpx; }", widget.GetCssName(), padding))
  return widget
 }
}

func (f Box) PaddingBottom(padding int) Box {
 return func() *gtk.Box {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-bottom: %dpx; }", widget.GetCssName(), padding))
  return widget
 }
}

func (f Box) PaddingEnd(padding int) Box {
 return func() *gtk.Box {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-right: %dpx; }", widget.GetCssName(), padding))
  return widget
 }
}

func (f Box) PaddingStart(padding int) Box {
 return func() *gtk.Box {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-left: %dpx; }", widget.GetCssName(), padding))
  return widget
 }
}

func (f Box) PaddingTop(padding int) Box {
 return func() *gtk.Box {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-top: %dpx; }", widget.GetCssName(), padding))
  return widget
 }
}

func (f Box) ToGTK() *gtk.Widget {
	val := f()
	return &val.Widget
}

func (f Box) VAlign(align gtk.Align) Box {
 return func() *gtk.Box {
  widget := f()
  widget.SetValign(align)
  return widget
 }
}

func (f Box) VExpand(expand bool) Box {
 return func() *gtk.Box {
  widget := f()
  widget.SetVexpand(expand)
  return widget
 }
}

func (f Box) Visible(visible bool) Box {
 return func() *gtk.Box {
  widget := f()
  widget.SetVisible(visible)
  return widget
 }
}

func (f Box) VMargin(vertical int) Box {
 return func() *gtk.Box {
  widget := f()
  widget.SetMarginTop(vertical)
  widget.SetMarginBottom(vertical)
  return widget
 }
}

