package schwifty

import (
	"fmt"
	"log/slog"
	"reflect"

	"codeberg.org/dergs/tidalwave/internal/g"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var logger = slog.With("library", "schwifty")

func ResolveWidget(value any) *gtk.Widget {
	t := reflect.TypeOf(value)

	if t.AssignableTo(reflect.TypeFor[*gtk.Widget]()) {
		logger.Debug("resolved widget from *gtk.Widget")
		fmt.Println("Was Widget", value.(*gtk.Widget).GoPointer())
		return value.(*gtk.Widget)
	}

	if t.Kind() == reflect.Pointer && t.Elem().Kind() == reflect.Struct {
		if field := reflect.ValueOf(value).Elem().FieldByName("Widget"); field.IsValid() {
			logger.Debug("resolved widget from specialized *gtk.Widget")
			widget := field.Interface().(gtk.Widget)
			return &widget
		}
	}

	if t.AssignableTo(reflect.TypeFor[BaseWidgetable]()) {
		logger.Debug("resolved widget from Widgetable")
		return value.(BaseWidgetable).ToGTK()
	}

	logger.Warn("failed to resolve widget")
	return nil
}

func UnrefOnDestroy(onWidget *gtk.Widget, widget *gtk.Widget) *func(gtk.Widget) {
	var id uint32
	cb := g.Ptr(func(w gtk.Widget) {
		widget.Unref()
		widget = nil
		w.DisconnectSignal(id)
	})
	id = onWidget.ConnectDestroy(cb)
	return cb
}
