package router

import (
	"log/slog"
	"time"

	"codeberg.org/dergs/tidalwave/internal/g"
)

var logger = slog.With("module", "router")

func Navigate(path string, params Params) {
	navigate(path, params, false)
}

func navigate(path string, params Params, offRecord bool) {
	if history.IsCurrentlyOn(path, params) && !offRecord {
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
		entry := &HistoryEntry{Path: path, Params: params, PageTitle: response.PageTitle, ExpiresAt: response.ExpiresAt}
		if response.Toolbar != nil {
			entry.Toolbar = response.Toolbar.ToGTK()
		}
		if response.View != nil {
			entry.View = response.View.ToGTK()
		}
		if !shouldCache {
			entry.ExpiresAt = g.Ptr(time.Now())
		}
		if !offRecord {
			history.Push(entry)
		}
		handleNavigationComplete(entry)

	}(path, params, handler)
}

func Back() {
	if len(history.Entries) == 0 {
		return
	}

	previous := history.Pop()
	if previous == nil {
		return
	}

	OnNavigate.Notify(previous.Path)
	if previous.View != nil {
		handleNavigationComplete(previous)
	} else {
		navigate(previous.Path, previous.Params, true)
	}
}

func Refresh() {
	if history.Current == nil {
		return
	}

	history.Current.Toolbar = nil
	history.Current.View = nil

	navigate(history.Current.Path, history.Current.Params, true)
}
