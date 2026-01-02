package schwifty

import (
	"log/slog"
	"reflect"

	"codeberg.org/dergs/tidalwave/pkg/schwifty/callback"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/tracking"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var logger = slog.With("library", "schwifty")

var (
	shouldLogLifecycle = false
)

func OnMainThread(cb callback.MainThreadCallback, param uintptr) uint {
	return callback.OnMainThread(cb, param)
}

func OnMainThreadOnce(cb func(u uintptr), param uintptr) uint {
	return callback.OnMainThreadOnce(cb, param)
}

func ResolveWidget(value any) *gtk.Widget {
	t := reflect.TypeOf(value)

	if value == nil {
		return nil
	}

	if t.AssignableTo(reflect.TypeFor[*gtk.Widget]()) {
		if shouldLogLifecycle {
			logger.Debug("resolved widget from *gtk.Widget")
		}
		return value.(*gtk.Widget)
	}

	if t.Kind() == reflect.Pointer && t.Elem().Kind() == reflect.Struct {
		if field := reflect.ValueOf(value).Elem().FieldByName("Widget"); field.IsValid() {
			if shouldLogLifecycle {
				logger.Debug("resolved widget from specialized *gtk.Widget")
			}
			widget := field.Interface().(gtk.Widget)
			return &widget
		}
	}

	if t.AssignableTo(reflect.TypeFor[BaseWidgetable]()) {
		if shouldLogLifecycle {
			logger.Debug("resolved widget from Widgetable")
		}
		return value.(BaseWidgetable).ToGTK()
	}

	logger.Warn("failed to resolve widget")
	return nil
}

func ResolvePopover(value any) *gtk.Popover {
	t := reflect.TypeOf(value)

	if t.AssignableTo(reflect.TypeFor[*gtk.Popover]()) {
		if shouldLogLifecycle {
			logger.Debug("resolved popover from *gtk.Popover")
		}
		return value.(*gtk.Popover)
	}

	if t.AssignableTo(reflect.TypeFor[Popover]()) {
		if shouldLogLifecycle {
			logger.Debug("resolved popover from Popover")
		}
		return value.(Popover)()
	}

	logger.Warn("failed to resolve widget")
	return nil
}

func SetLogLifecycle(enabled bool) {
	shouldLogLifecycle = enabled
	callback.SetLogLifecycle(enabled)
	tracking.SetLogLifecycle(enabled)
}
