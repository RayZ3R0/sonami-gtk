package callback

import (
	"log/slog"
	"reflect"
	"sync"

	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var (
	logger              = slog.With("library", "schwifty", "module", "callback")
	widgetCallbacks     = make(map[uintptr]map[string][]any)
	widgetCallbacksLock = sync.RWMutex{}
)

func CallbackHandler[T any](widget gtk.Widget, signal string, args ...any) []T {
	widget.Ref()
	defer widget.Unref()
	widgetCallbacksLock.Lock()
	defer widgetCallbacksLock.Unlock()

	// Check if the widget has any callbacks registered
	if allCallbacks, ok := widgetCallbacks[widget.GoPointer()]; !ok {
		return nil
	} else if signalCallbacks, ok := allCallbacks[signal]; !ok {
		return nil
	} else if len(signalCallbacks) >= 0 {
		returnValues := make([]T, len(signalCallbacks))
		for i, callback := range signalCallbacks {
			reflectArgs := make([]reflect.Value, len(args))
			for i, arg := range args {
				reflectArgs[i] = reflect.ValueOf(arg)
			}
			result := reflect.ValueOf(callback).Call(reflectArgs)
			if len(result) > 0 {
				returnValues[i] = result[0].Interface().(T)
			}
		}
		logger.Debug("executed callback", "ptr", widget.GoPointer(), "signal", signal, "handlers", len(signalCallbacks))
		return returnValues
	}
	return nil
}

func HandleCallback(widget gtk.Widget, signal string, callback any) {
	widget.Ref()
	defer widget.Unref()
	widgetCallbacksLock.Lock()
	defer widgetCallbacksLock.Unlock()

	id := widget.GoPointer()

	// Check if the widget has any callbacks registered
	allCallbacks, ok := widgetCallbacks[id]
	if !ok {
		allCallbacks = make(map[string][]any)
	}

	// Check if the signal has any callbacks registered
	signalCallbacks, ok := allCallbacks[signal]
	if !ok {
		signalCallbacks = make([]any, 0)
	}

	// Add the callback to the list of callbacks for the signal
	signalCallbacks = append(signalCallbacks, callback)
	allCallbacks[signal] = signalCallbacks
	widgetCallbacks[id] = allCallbacks
	logger.Debug("registered callback", "ptr", widget.GoPointer(), "signal", signal)
}

func DeleteCallbacks(widget gtk.Widget) {
	widget.Ref()
	defer widget.Unref()
	widgetCallbacksLock.Lock()
	defer widgetCallbacksLock.Unlock()
	delete(widgetCallbacks, widget.GoPointer())
	logger.Debug("deleted all callbacks", "ptr", widget.GoPointer())
}
