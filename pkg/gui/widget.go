package gui

import (
	"fmt"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotkit/gtkutil/cssutil"
)

type WidgetImpl[T any] struct {
	gtk.Widgetter
	widget gtk.Widget
	real   T
}

func (w *WidgetImpl[T]) AddCSSClass(class string) T {
	w.widget.AddCSSClass(class)
	return w.real
}

func (w *WidgetImpl[T]) Background(color string) T {
	cssutil.Apply(w, fmt.Sprintf("%s { background-color: %s; }", w.widget.CSSName(), color))
	return w.real
}

func (w *WidgetImpl[T]) CallCSSApplier(applier func(gtk.Widgetter)) T {
	applier(w)
	return w.real
}

func (w *WidgetImpl[T]) CornerRadius(radius int) T {
	cssutil.Apply(w, fmt.Sprintf(w.widget.CSSName()+` { border-radius: %dpx; }`, radius))
	return w.real
}

func (w *WidgetImpl[T]) CSS(css string) T {
	cssutil.Apply(w, css)
	return w.real
}

func (w *WidgetImpl[T]) Focusable(focusable bool) T {
	w.widget.SetFocusable(focusable)
	return w.real
}

func (w *WidgetImpl[T]) FocusOnClick(focusOnClick bool) T {
	w.widget.SetFocusOnClick(focusOnClick)
	return w.real
}

func (w *WidgetImpl[T]) HAlign(align gtk.Align) T {
	w.widget.SetHAlign(align)
	return w.real
}

func (w *WidgetImpl[T]) HExpand(expand bool) T {
	w.widget.SetHExpand(expand)
	return w.real
}

func (w *WidgetImpl[T]) HMargin(horizontal int) T {
	w.widget.SetMarginStart(horizontal)
	w.widget.SetMarginEnd(horizontal)
	return w.real
}

func (w *WidgetImpl[T]) Margin(margin int) T {
	w.widget.SetMarginStart(margin)
	w.widget.SetMarginEnd(margin)
	w.widget.SetMarginTop(margin)
	w.widget.SetMarginBottom(margin)
	return w.real
}

func (w *WidgetImpl[T]) MarginLeft(start int) T {
	w.widget.SetMarginStart(start)
	return w.real
}

func (w *WidgetImpl[T]) MarginRight(end int) T {
	w.widget.SetMarginEnd(end)
	return w.real
}

func (w *WidgetImpl[T]) MarginTop(top int) T {
	w.widget.SetMarginTop(top)
	return w.real
}

func (w *WidgetImpl[T]) MarginBottom(bottom int) T {
	w.widget.SetMarginBottom(bottom)
	return w.real
}

func (w *WidgetImpl[T]) Opacity(opacity float64) T {
	w.widget.SetOpacity(opacity)
	return w.real
}

func (w *WidgetImpl[T]) Overflow(overflow gtk.Overflow) T {
	w.widget.SetOverflow(overflow)
	return w.real
}

func (w *WidgetImpl[T]) Padding(paddingPx int) T {
	cssutil.Apply(w, fmt.Sprintf("%s { padding: %dpx; }", w.widget.CSSName(), paddingPx))
	return w.real
}

func (w *WidgetImpl[T]) PaddingStart(paddingPx int) T {
	cssutil.Apply(w, fmt.Sprintf("%s { padding-left: %dpx; }", w.widget.CSSName(), paddingPx))
	return w.real
}

func (w *WidgetImpl[T]) PaddingEnd(paddingPx int) T {
	cssutil.Apply(w, fmt.Sprintf("%s { padding-right: %dpx; }", w.widget.CSSName(), paddingPx))
	return w.real
}

func (w *WidgetImpl[T]) PaddingTop(paddingPx int) T {
	cssutil.Apply(w, fmt.Sprintf("%s { padding-top: %dpx; }", w.widget.CSSName(), paddingPx))
	return w.real
}

func (w *WidgetImpl[T]) PaddingBottom(paddingPx int) T {
	cssutil.Apply(w, fmt.Sprintf("%s { padding-bottom: %dpx; }", w.widget.CSSName(), paddingPx))
	return w.real
}

func (w *WidgetImpl[T]) RemoveCSSClass(className string) T {
	w.widget.RemoveCSSClass(className)
	return w.real
}

func (w *WidgetImpl[T]) VAlign(align gtk.Align) T {
	w.widget.SetVAlign(align)
	return w.real
}

func (w *WidgetImpl[T]) VExpand(expand bool) T {
	w.widget.SetVExpand(expand)
	return w.real
}

func (w *WidgetImpl[T]) VMargin(vertical int) T {
	w.widget.SetMarginTop(vertical)
	w.widget.SetMarginBottom(vertical)
	return w.real
}

func (w *WidgetImpl[T]) GTKWidget() *gtk.Widget {
	return &w.widget
}
