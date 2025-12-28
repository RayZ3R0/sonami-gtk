package schwifty

import (
	"fmt"

	"codeberg.org/dergs/tidalwave/pkg/schwifty/css"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

type Widget func() *WrappedWidget

func (f Widget) AddController(controller *gtk.EventController) Widget {
 return func() *WrappedWidget {
  widget := f()
  widget.AddController(controller)
  return widget
 }
}

func (f Widget) Background(color string) Widget {
 return func() *WrappedWidget {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { background-color: %s; }", widget.GetCssName(), color))
  return widget
 }
}

func (f Widget) CornerRadius(radius int) Widget {
 return func() *WrappedWidget {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { border-radius: %dpx; }", widget.GetCssName(), radius))
  return widget
 }
}

func (f Widget) CSS(css string) Widget {
	return func() *WrappedWidget {
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

func (f Widget) Focusable(focusable bool) Widget {
	return func() *WrappedWidget {
		widget := f()
		widget.SetFocusable(focusable)
		return widget
	}
}

func (f Widget) FocusOnClick(focusOnClick bool) Widget {
	return func() *WrappedWidget {
		widget := f()
		widget.SetFocusOnClick(focusOnClick)
		return widget
	}
}

func (f Widget) HAlign(align gtk.Align) Widget {
	return func() *WrappedWidget {
		widget := f()
		widget.SetHalign(align)
		return widget
	}
}

func (f Widget) HExpand(expand bool) Widget {
	return func() *WrappedWidget {
		widget := f()
		widget.SetHexpand(expand)
		return widget
	}
}

func (f Widget) HMargin(horizontal int) Widget {
	return func() *WrappedWidget {
		widget := f()
		widget.SetMarginEnd(horizontal)
		widget.SetMarginStart(horizontal)
		return widget
	}
}

func (f Widget) Margin(margin int) Widget {
	return func() *WrappedWidget {
		widget := f()
		widget.SetMarginBottom(margin)
		widget.SetMarginEnd(margin)
		widget.SetMarginStart(margin)
		widget.SetMarginTop(margin)
		return widget
	}
}

func (f Widget) MarginBottom(bottom int) Widget {
	return func() *WrappedWidget {
		widget := f()
		widget.SetMarginBottom(bottom)
		return widget
	}
}

func (f Widget) MarginEnd(end int) Widget {
	return func() *WrappedWidget {
		widget := f()
		widget.SetMarginEnd(end)
		return widget
	}
}

func (f Widget) MarginStart(start int) Widget {
	return func() *WrappedWidget {
		widget := f()
		widget.SetMarginStart(start)
		return widget
	}
}

func (f Widget) MarginTop(top int) Widget {
	return func() *WrappedWidget {
		widget := f()
		widget.SetMarginTop(top)
		return widget
	}
}

func (f Widget) Opacity(opacity float64) Widget {
 return func() *WrappedWidget {
  widget := f()
  widget.SetOpacity(opacity)
  return widget
 }
}

func (f Widget) Overflow(overflow gtk.Overflow) Widget {
 return func() *WrappedWidget {
  widget := f()
  widget.SetOverflow(overflow)
  return widget
 }
}

func (f Widget) ToGTK() *gtk.Widget {
	val := f()
	return &val.Widget
}

func (f Widget) VAlign(align gtk.Align) Widget {
 return func() *WrappedWidget {
  widget := f()
  widget.SetValign(align)
  return widget
 }
}

func (f Widget) VExpand(expand bool) Widget {
 return func() *WrappedWidget {
  widget := f()
  widget.SetVexpand(expand)
  return widget
 }
}

func (f Widget) Visible(visible bool) Widget {
 return func() *WrappedWidget {
  widget := f()
  widget.SetVisible(visible)
  return widget
 }
}

func (f Widget) VMargin(vertical int) Widget {
 return func() *WrappedWidget {
  widget := f()
  widget.SetMarginTop(vertical)
  widget.SetMarginBottom(vertical)
  return widget
 }
}

