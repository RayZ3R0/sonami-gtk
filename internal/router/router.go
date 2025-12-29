package router

import (
	"log/slog"
	"time"
)

var logger = slog.With("module", "router")

func Navigate(path string, params Params) {
	if history.IsCurrentlyOn(path, params) {
		logger.Debug("skipped navigation as we are already on the same page")
		return
	}

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
		entry := HistoryEntry{Path: path, Params: params, PageTitle: response.PageTitle}
		if response.Toolbar != nil {
			entry.Toolbar = response.Toolbar.ToGTK()
		}
		if response.View != nil {
			entry.View = response.View.ToGTK()
		}
		if shouldCache {
			entry.Response = response
		}
		history.Push(entry)
		handleNavigationComplete(entry)

	}(path, params, handler)
}

func Back() {
	if len(history.array) < 1 {
		return
	}

	current := history.Current()
	previous := history.Pop()
	if previous == nil {
		return
	}

	// Ensure we are unmapping whatever we were showing previously
	defer func(entry HistoryEntry) {
		if entry.Toolbar != nil {
			entry.Toolbar.Unref()
		}
		if entry.View != nil {
			entry.View.Unref()
		}
	}(*current)

	OnNavigate.Notify(previous.Path)
	if previous.Response != nil {
		handleNavigationComplete(*previous)
	} else {
		Navigate(previous.Path, previous.Params)
	}
}
