package gui

import (
	"reflect"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type WrapperImpl struct {
	*WidgetImpl[*WrapperImpl]
	widget gtk.Widgetter
}

var Wrapper = func(root gtk.Widgetter) *WrapperImpl {
	if root == nil {
		panic("root widget cannot be nil")
	}
	widget := reflect.ValueOf(root).Elem().FieldByName("Widget").Interface()
	impl := &WrapperImpl{nil, root}
	impl.WidgetImpl = &WidgetImpl[*WrapperImpl]{root, widget.(gtk.Widget), impl}
	return impl
}
