package main

import (
	"fmt"
	"os"
	"strings"
)

var glue = `package schwifty

import (
	"fmt"

	"codeberg.org/dergs/tidalwave/pkg/schwifty/css"
	"github.com/jwijenbergh/puregotk/v4/gtk"{{IMPORTS}}
)

type {{TYPE}} func() {{BASE_TYPE}}

func (f {{TYPE}}) AddController(controller *gtk.EventController) {{TYPE}} {
 return func() {{BASE_TYPE}} {
  widget := f()
  widget.AddController(controller)
  return widget
 }
}

func (f {{TYPE}}) Background(color string) {{TYPE}} {
 return func() {{BASE_TYPE}} {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { background-color: %s; }", widget.GetCssName(), color))
  return widget
 }
}

func (f {{TYPE}}) CornerRadius(radius int) {{TYPE}} {
 return func() {{BASE_TYPE}} {
  widget := f()
  css.Apply(&widget.Widget, fmt.Sprintf("%s { border-radius: %dpx; }", widget.GetCssName(), radius))
  return widget
 }
}

func (f {{TYPE}}) CSS(css string) {{TYPE}} {
	return func() {{BASE_TYPE}} {
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

func (f {{TYPE}}) Focusable(focusable bool) {{TYPE}} {
	return func() {{BASE_TYPE}} {
		widget := f()
		widget.SetFocusable(focusable)
		return widget
	}
}

func (f {{TYPE}}) FocusOnClick(focusOnClick bool) {{TYPE}} {
	return func() {{BASE_TYPE}} {
		widget := f()
		widget.SetFocusOnClick(focusOnClick)
		return widget
	}
}

func (f {{TYPE}}) HAlign(align gtk.Align) {{TYPE}} {
	return func() {{BASE_TYPE}} {
		widget := f()
		widget.SetHalign(align)
		return widget
	}
}

func (f {{TYPE}}) HExpand(expand bool) {{TYPE}} {
	return func() {{BASE_TYPE}} {
		widget := f()
		widget.SetHexpand(expand)
		return widget
	}
}

func (f {{TYPE}}) HMargin(horizontal int) {{TYPE}} {
	return func() {{BASE_TYPE}} {
		widget := f()
		widget.SetMarginEnd(horizontal)
		widget.SetMarginStart(horizontal)
		return widget
	}
}

func (f {{TYPE}}) Margin(margin int) {{TYPE}} {
	return func() {{BASE_TYPE}} {
		widget := f()
		widget.SetMarginBottom(margin)
		widget.SetMarginEnd(margin)
		widget.SetMarginStart(margin)
		widget.SetMarginTop(margin)
		return widget
	}
}

func (f {{TYPE}}) MarginBottom(bottom int) {{TYPE}} {
	return func() {{BASE_TYPE}} {
		widget := f()
		widget.SetMarginBottom(bottom)
		return widget
	}
}

func (f {{TYPE}}) MarginEnd(end int) {{TYPE}} {
	return func() {{BASE_TYPE}} {
		widget := f()
		widget.SetMarginEnd(end)
		return widget
	}
}

func (f {{TYPE}}) MarginStart(start int) {{TYPE}} {
	return func() {{BASE_TYPE}} {
		widget := f()
		widget.SetMarginStart(start)
		return widget
	}
}

func (f {{TYPE}}) MarginTop(top int) {{TYPE}} {
	return func() {{BASE_TYPE}} {
		widget := f()
		widget.SetMarginTop(top)
		return widget
	}
}

func (f {{TYPE}}) Opacity(opacity float64) {{TYPE}} {
 return func() {{BASE_TYPE}} {
  widget := f()
  widget.SetOpacity(opacity)
  return widget
 }
}

func (f {{TYPE}}) Overflow(overflow gtk.Overflow) {{TYPE}} {
 return func() {{BASE_TYPE}} {
  widget := f()
  widget.SetOverflow(overflow)
  return widget
 }
}

func (f {{TYPE}}) ToGTK() *gtk.Widget {
	val := f()
	return &val.Widget
}

func (f {{TYPE}}) VAlign(align gtk.Align) {{TYPE}} {
 return func() {{BASE_TYPE}} {
  widget := f()
  widget.SetValign(align)
  return widget
 }
}

func (f {{TYPE}}) VExpand(expand bool) {{TYPE}} {
 return func() {{BASE_TYPE}} {
  widget := f()
  widget.SetVexpand(expand)
  return widget
 }
}

func (f {{TYPE}}) Visible(visible bool) {{TYPE}} {
 return func() {{BASE_TYPE}} {
  widget := f()
  widget.SetVisible(visible)
  return widget
 }
}

func (f {{TYPE}}) VMargin(vertical int) {{TYPE}} {
 return func() {{BASE_TYPE}} {
  widget := f()
  widget.SetMarginTop(vertical)
  widget.SetMarginBottom(vertical)
  return widget
 }
}

`

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: gen <type> <base_type>")
		os.Exit(1)
	}

	typ := os.Args[1]
	baseType := os.Args[2]

	imports := ""
	if strings.Contains(baseType, "adw.") {
		imports += "\n\t\"github.com/jwijenbergh/puregotk/v4/adw\""
	}

	output := strings.ReplaceAll(glue, "{{TYPE}}", typ)
	output = strings.ReplaceAll(output, "{{BASE_TYPE}}", baseType)
	output = strings.ReplaceAll(output, "{{IMPORTS}}", imports)
	os.WriteFile(strings.ToLower(typ)+"_generated.go", []byte(output), 0644)
}
