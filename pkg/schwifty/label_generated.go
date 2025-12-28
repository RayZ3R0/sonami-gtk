package schwifty

import (
	"fmt"

	"codeberg.org/dergs/tidalwave/pkg/schwifty/css"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

type Label func() *gtk.Label

func (f Label) AddController(controller *gtk.EventController) Label {
 return func() *gtk.Label {
  widget := f()
  widget.AddController(controller)
  return widget
 }
}

func (f Label) Background(color string) Label {
 return func() *gtk.Label {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { background-color: %s; }", widget.GetCssName(), color))
  return widget
 }
}

func (f Label) CornerRadius(radius int) Label {
 return func() *gtk.Label {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { border-radius: %dpx; }", widget.GetCssName(), radius))
  return widget
 }
}

func (f Label) CSS(css string) Label {
	return func() *gtk.Label {
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

func (f Label) Focusable(focusable bool) Label {
	return func() *gtk.Label {
		widget := f()
		widget.SetFocusable(focusable)
		return widget
	}
}

func (f Label) FocusOnClick(focusOnClick bool) Label {
	return func() *gtk.Label {
		widget := f()
		widget.SetFocusOnClick(focusOnClick)
		return widget
	}
}

func (f Label) HAlign(align gtk.Align) Label {
	return func() *gtk.Label {
		widget := f()
		widget.SetHalign(align)
		return widget
	}
}

func (f Label) HExpand(expand bool) Label {
	return func() *gtk.Label {
		widget := f()
		widget.SetHexpand(expand)
		return widget
	}
}

func (f Label) HMargin(horizontal int) Label {
	return func() *gtk.Label {
		widget := f()
		widget.SetMarginEnd(horizontal)
		widget.SetMarginStart(horizontal)
		return widget
	}
}

func (f Label) Margin(margin int) Label {
	return func() *gtk.Label {
		widget := f()
		widget.SetMarginBottom(margin)
		widget.SetMarginEnd(margin)
		widget.SetMarginStart(margin)
		widget.SetMarginTop(margin)
		return widget
	}
}

func (f Label) MarginBottom(bottom int) Label {
	return func() *gtk.Label {
		widget := f()
		widget.SetMarginBottom(bottom)
		return widget
	}
}

func (f Label) MarginEnd(end int) Label {
	return func() *gtk.Label {
		widget := f()
		widget.SetMarginEnd(end)
		return widget
	}
}

func (f Label) MarginStart(start int) Label {
	return func() *gtk.Label {
		widget := f()
		widget.SetMarginStart(start)
		return widget
	}
}

func (f Label) MarginTop(top int) Label {
	return func() *gtk.Label {
		widget := f()
		widget.SetMarginTop(top)
		return widget
	}
}

func (f Label) Opacity(opacity float64) Label {
 return func() *gtk.Label {
  widget := f()
  widget.SetOpacity(opacity)
  return widget
 }
}

func (f Label) Overflow(overflow gtk.Overflow) Label {
 return func() *gtk.Label {
  widget := f()
  widget.SetOverflow(overflow)
  return widget
 }
}

func (f Label) ToGTK() *gtk.Widget {
	val := f()
	return &val.Widget
}

func (f Label) VAlign(align gtk.Align) Label {
 return func() *gtk.Label {
  widget := f()
  widget.SetValign(align)
  return widget
 }
}

func (f Label) VExpand(expand bool) Label {
 return func() *gtk.Label {
  widget := f()
  widget.SetVexpand(expand)
  return widget
 }
}

func (f Label) Visible(visible bool) Label {
 return func() *gtk.Label {
  widget := f()
  widget.SetVisible(visible)
  return widget
 }
}

func (f Label) VMargin(vertical int) Label {
 return func() *gtk.Label {
  widget := f()
  widget.SetMarginTop(vertical)
  widget.SetMarginBottom(vertical)
  return widget
 }
}

