package gtk

import (
	"reflect"

	"codeberg.org/dergs/tonearm/pkg/schwifty/callback"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

//go:generate go run codeberg.org/dergs/tonearm/pkg/schwifty/gen Widget *WrappedWidget gtk
type WrappedWidget struct {
	gtk.Widget
}

type BaseWidgetable interface {
	ToGTK() *gtk.Widget
}

type Widgetable[T any] interface {
	BaseWidgetable
	AddController(controller *gtk.EventController) T
	CSS(css string) T
	Focusable(focusable bool) T
	FocusOnClick(focusOnClick bool) T
	HAlign(align gtk.Align) T
	HExpand(expand bool) T
	HMargin(horizontal int) T
	Margin(margin int) T
	MarginBottom(bottom int) T
	MarginEnd(end int) T
	MarginStart(start int) T
	MarginTop(top int) T
	Opacity(opacity float64) T
	Overflow(overflow gtk.Overflow) T
	VAlign(align gtk.Align) T
	VExpand(expand bool) T
	Visible(visible bool) T
	VMargin(vertical int) T
}

func ResolveWidget(value any) *gtk.Widget {
	t := reflect.TypeOf(value)

	if value == nil {
		return nil
	}

	if t.AssignableTo(reflect.TypeFor[*gtk.Widget]()) {
		// if shouldLogLifecycle {
		// 	logger.Debug("resolved widget from *gtk.Widget")
		// }
		return value.(*gtk.Widget)
	}

	if t.AssignableTo(reflect.TypeFor[BaseWidgetable]()) {
		// if shouldLogLifecycle {
		// 	logger.Debug("resolved widget from Widgetable")
		// }
		return value.(BaseWidgetable).ToGTK()
	}

	if t.Kind() == reflect.Pointer && t.Elem().Kind() == reflect.Struct {
		if field := reflect.ValueOf(value).Elem().FieldByName("Widget"); field.IsValid() {
			// if shouldLogLifecycle {
			// 	logger.Debug("resolved widget from specialized *gtk.Widget")
			// }
			widget := field.Interface().(gtk.Widget)
			return &widget
		}
	}

	// logger.Warn("failed to resolve widget")
	return nil
}

func ResolveWidgetOnMain(value any) (channel chan *gtk.Widget) {
	t := reflect.TypeOf(value)
	channel = make(chan *gtk.Widget, 1)

	if value == nil {
		channel <- nil
		return
	}

	if t.AssignableTo(reflect.TypeFor[*gtk.Widget]()) {
		channel <- value.(*gtk.Widget)
		return
	}

	if t.AssignableTo(reflect.TypeFor[BaseWidgetable]()) {
		callback.OnMainThreadOncePure(func() {
			channel <- value.(BaseWidgetable).ToGTK()
		})
		return channel
	}

	if t.Kind() == reflect.Pointer && t.Elem().Kind() == reflect.Struct {
		if field := reflect.ValueOf(value).Elem().FieldByName("Widget"); field.IsValid() {
			widget := field.Interface().(gtk.Widget)
			channel <- &widget
			return
		}
	}

	channel <- nil
	return
}
