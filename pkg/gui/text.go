package gui

import (
	"fmt"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotk4/pkg/pango"
	"github.com/diamondburned/gotkit/gtkutil/cssutil"
)

type TextFunc func(text string) *TextImpl

type TextImpl struct {
	*WidgetImpl[*TextImpl]
	label *gtk.Label
}

var Text = func(text string) *TextImpl {
	label := gtk.NewLabel(text)
	impl := &TextImpl{
		WidgetImpl: nil,
		label:      label,
	}
	impl.WidgetImpl = &WidgetImpl[*TextImpl]{label, label.Widget, impl}
	return impl
}

func (t *TextImpl) Ellipsis(mode pango.EllipsizeMode) *TextImpl {
	t.label.SetEllipsize(mode)
	return t
}

func (t *TextImpl) FontSize(weight int) *TextImpl {
	cssutil.Apply(t, fmt.Sprintf("label { font-size: %dpx; }", weight))
	return t
}

func (t *TextImpl) FontWeight(weight int) *TextImpl {
	cssutil.Apply(t, fmt.Sprintf("label { font-weight: %d; }", weight))
	return t
}

func (t *TextImpl) GTKWidget() *gtk.Label {
	return t.label
}

func (t *TextImpl) Text(text string) *TextImpl {
	t.label.SetText(text)
	return t
}
