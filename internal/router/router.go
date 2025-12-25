package router

import (
	"log/slog"
	"time"
)

var logger = slog.With("module", "router")

func Navigate(path string, params Params) {
	logger.Debug("navigation started")
	// We are starting to navigate, notify the rest of the application
	OnNavigate.Notify(path)

	// If the route is not registered, use the not found handler
	handler, ok := routeMap[path]
	if !ok {
		logger.Info("no handler found", "path", path)
		handler = notFoundHandler
	}

	startTime := time.Now()
	go func(path string, params Params, handler Handler) {
		response, shouldCache := executeHandler(handler, params)
		logger.Info("navigation completed", "path", path, "params", params, "duration_ms", time.Since(startTime).Milliseconds(), "should_cache", shouldCache)
		if shouldCache {
			history.Push(HistoryEntry{Path: path, Params: params, Response: response})
		} else {
			history.Push(HistoryEntry{Path: path, Params: params})
		}
		handleResponse(path, params, response)

	}(path, params, handler)
}

func Back() {
	if len(history.array) < 1 {
		return
	}

	previous := history.Pop()
	if previous == nil {
		return
	}

	if previous.Response != nil {
		handleResponse(previous.Path, previous.Params, previous.Response)
	} else {
		Navigate(previous.Path, previous.Params)
	}
}
