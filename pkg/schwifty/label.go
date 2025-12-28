package schwifty

import (
	"fmt"

	"codeberg.org/dergs/tidalwave/pkg/schwifty/css"
	"github.com/jwijenbergh/puregotk/v4/gtk"
	"github.com/jwijenbergh/puregotk/v4/pango"
)

//go:generate go run codeberg.org/dergs/tidalwave/pkg/schwifty/gen Label *gtk.Label

// Sets the mode used to ellipsize the text.
//
// The text will be ellipsized if there is not
// enough space to render the entire string.
func (f Label) Ellipsis(ellipsis pango.EllipsizeMode) Label {
	return func() *gtk.Label {
		label := f()
		label.SetEllipsize(ellipsis)
		return label
	}
}

// Sets the font size for the label text.
//
// The size parameter determines the size of the text in points.
// Common values include 12 for small text and 16 for larger text,
// though the exact range and supported values depend on the font family.
//
// This function modifies the font attributes of the label to apply
// the specified size to the displayed text.
func (f Label) FontSize(size int) Label {
	return func() *gtk.Label {
		label := f()
		css.Apply(&label.Widget, fmt.Sprintf("label { font-size: %dpx; }", size))
		return label
	}
}

// Sets the font weight for the label text.
//
// The weight parameter determines how bold or light the text appears.
// Common values include 400 for normal weight and 700 for bold weight,
// though the exact range and supported values depend on the font family.
//
// This function modifies the font attributes of the label to apply
// the specified weight to the displayed text.
func (f Label) FontWeight(weight int) Label {
	return func() *gtk.Label {
		label := f()
		css.Apply(&label.Widget, fmt.Sprintf("label { font-weight: %d; }", weight))
		return label
	}
}

// Sets the alignment of lines in the label relative to each other.
//
// This function has no effect on labels containing only a single line.
//
// [enum@Gtk.Justification.left] is the default value when the widget
// is first created with [ctor@Gtk.Label.new].
//
// If you instead want to set the alignment of the label as a whole,
// use [method@Gtk.Widget.set_halign] instead.
func (f Label) Justify(justify gtk.Justification) Label {
	return func() *gtk.Label {
		label := f()
		label.SetJustify(justify)
		return label
	}
}

// Sets the text for the label.
//
// It overwrites any text that was there before and clears any
// previously set mnemonic accelerators, and sets the
// [property@Gtk.Label:use-underline] and
// [property@Gtk.Label:use-markup] properties to false.
//
// Also see [method@Gtk.Label.set_markup].
func (f Label) Text(text string) Label {
	return func() *gtk.Label {
		label := f()
		label.SetText(text)
		return label
	}
}

// Toggles line wrapping within the label.
//
// True makes it break lines if text exceeds the widget’s size.
// false lets the text get cut off by the edge of the widget if
// it exceeds the widget size.
//
// Note that setting line wrapping to true does not make the label
// wrap at its parent widget’s width, because GTK widgets conceptually
// can’t make their requisition depend on the parent  widget’s size.
// For a label that wraps at a specific position, set the label’s width
// using [method@Gtk.Widget.set_size_request].
func (f Label) Wrap(wrap bool) Label {
	return func() *gtk.Label {
		label := f()
		label.SetWrap(wrap)
		return label
	}
}
