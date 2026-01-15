package router

import "codeberg.org/dergs/tonearm/internal/signals"

// NavigationStarted is a signal that is emitted when the router starts navigating to a new path.
// The path parameter is the new path that the router is navigating to.
var NavigationStarted = signals.NewStatelessSignal[string]()

// NavigationCompleted is a signal that is emitted when the router completes a navigation.
// The response produced by the page handler cannot be nil.
var NavigationCompleted = signals.NewStatelessSignal[HistoryEntry]()

// HistoryUpdated is a signal that is emitted when the router updates the history.
// The history parameter is the new history that the router is updating.
// The history parameter cannot be nil.
var HistoryUpdated = signals.NewStatelessSignal[*History]()
