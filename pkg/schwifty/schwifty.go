package schwifty

import (
	"log/slog"
	"reflect"

	"codeberg.org/dergs/tonearm/pkg/schwifty/callback"
	"codeberg.org/dergs/tonearm/pkg/schwifty/tracking"
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

func OnMainThreadOncePure(cb func()) uint {
	return OnMainThreadOnce(func(uintptr) { cb() }, 0)
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

	if t.AssignableTo(reflect.TypeFor[BaseWidgetable]()) {
		if shouldLogLifecycle {
			logger.Debug("resolved widget from Widgetable")
		}
		return value.(BaseWidgetable).ToGTK()
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

func ResolveTo[GtkType any, SchwiftyType BaseWidgetable](value any) (result GtkType) {
	t := reflect.TypeOf(value)

	if t.AssignableTo(reflect.TypeFor[GtkType]()) {
		if shouldLogLifecycle {
			logger.Debug("resolved generic from gtk")
		}
		result = value.(GtkType)
		return
	}

	if t.AssignableTo(reflect.TypeFor[SchwiftyType]()) {
		if shouldLogLifecycle {
			logger.Debug("resolved generic from schwifty")
		}
		result = reflect.ValueOf(value).Call([]reflect.Value{})[0].Interface().(GtkType)
		return
	}

	logger.Warn("failed to resolve generic")
	return
}

func SetLogLifecycle(enabled bool) {
	shouldLogLifecycle = enabled
	callback.SetLogLifecycle(enabled)
	tracking.SetLogLifecycle(enabled)
}
