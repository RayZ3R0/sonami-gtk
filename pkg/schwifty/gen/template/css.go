package schwifty

import (
	"fmt"

	"codeberg.org/dergs/tidalwave/pkg/schwifty/css"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func (f TEMPLATE_TYPE) Background(color string) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { background-color: %s; }", widget.GetCssName(), color))
		return widget
	}
}

func (f TEMPLATE_TYPE) CornerRadius(radius int) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { border-radius: %dpx; }", widget.GetCssName(), radius))
		return widget
	}
}

func (f TEMPLATE_TYPE) CSS(css string) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
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

func (f TEMPLATE_TYPE) HPadding(padding int) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-left: %dpx; padding-right: %dpx; }", widget.GetCssName(), padding, padding))
		return widget
	}
}

func (f TEMPLATE_TYPE) MinHeight(minHeight int) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { min-height: %dpx; }", widget.GetCssName(), minHeight))
		return widget
	}
}

func (f TEMPLATE_TYPE) MinWidth(minWidth int) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { min-width: %dpx; }", widget.GetCssName(), minWidth))
		return widget
	}
}

func (f TEMPLATE_TYPE) Padding(padding int) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f TEMPLATE_TYPE) PaddingBottom(padding int) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-bottom: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f TEMPLATE_TYPE) PaddingEnd(padding int) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-right: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f TEMPLATE_TYPE) PaddingStart(padding int) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-left: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f TEMPLATE_TYPE) PaddingTop(padding int) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-top: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f TEMPLATE_TYPE) VPadding(padding int) TEMPLATE_TYPE {
	return func() TEMPLATE_BASE_TYPE {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-bottom: %dpx; padding-top: %dpx; }", widget.GetCssName(), padding, padding))
		return widget
	}
}
