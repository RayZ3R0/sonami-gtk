package schwifty

import (
	"fmt"

	"codeberg.org/dergs/tidalwave/pkg/schwifty/css"
	"github.com/jwijenbergh/puregotk/v4/gtk"
	"github.com/jwijenbergh/puregotk/v4/adw"
)

type WindowTitle func() *adw.WindowTitle

func (f WindowTitle) AddController(controller *gtk.EventController) WindowTitle {
 return func() *adw.WindowTitle {
  widget := f()
  widget.AddController(controller)
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

