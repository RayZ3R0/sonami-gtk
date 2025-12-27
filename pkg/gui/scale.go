package gui

import "github.com/diamondburned/gotk4/pkg/gtk/v4"

type ScaleImpl struct {
	*WidgetImpl[*ScaleImpl]
	scale *gtk.Scale
}

func Scale(orientation gtk.Orientation) *ScaleImpl {
	scale := gtk.NewScale(orientation, nil)
	impl := &ScaleImpl{nil, scale}
	impl.WidgetImpl = &WidgetImpl[*ScaleImpl]{scale, scale.Widget, impl}
	return impl
}

func (s *ScaleImpl) DefaultValue(value float64) *ScaleImpl {
	s.scale.SetValue(value)
	return s
}

func (s *ScaleImpl) OnChange(callback func(gtk.ScrollType, float64) bool) *ScaleImpl {
	s.scale.ConnectChangeValue(callback)
	return s
}

func (s *ScaleImpl) GTKWidget() *gtk.Scale {
	return s.scale
}

func (s *ScaleImpl) Invert() *ScaleImpl {
	s.scale.SetInverted(!s.scale.Inverted())
	return s
}

func (s *ScaleImpl) Range(min, max float64) *ScaleImpl {
	s.scale.SetRange(min, max)
	return s
}
