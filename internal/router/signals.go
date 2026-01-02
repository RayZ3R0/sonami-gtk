package router

import "codeberg.org/dergs/tidalwave/internal/signals"

// OnNavigate is a signal that is emitted when the router starts navigating to a new path.
// The path parameter is the new path that the router is navigating to.
var OnNavigate = routerOnNavigateSignal{
	signals.NewSignal[func(path string) bool](),
}

type routerOnNavigateSignal struct {
	signals.Signal[func(path string) bool]
}

func (r *routerOnNavigateSignal) Notify(path string) {
	r.Signal.Notify(path)
}

// NavigationComplete is a signal that is emitted when the router completes a navigation.
// The response produced by the page handler cannot be nil.
var NavigationComplete = routerNavigationCompleteSignal{
	signals.NewSignal[func(entry HistoryEntry) bool](),
}

type routerNavigationCompleteSignal struct {
	signals.Signal[func(entry HistoryEntry) bool]
}

func (r *routerNavigationCompleteSignal) Notify(entry HistoryEntry) {
	r.Signal.Notify(entry)
}

// HistoryUpdated is a signal that is emitted when the router updates the history.
// The history parameter is the new history that the router is updating.
// The history parameter cannot be nil.
var HistoryUpdated = historyUpdatedSignal{
	signals.NewSignal[func(history *History) bool](),
}

type historyUpdatedSignal struct {
	signals.Signal[func(history *History) bool]
}

func (r *historyUpdatedSignal) Notify(history *History) {
	if history == nil {
		panic("history for HistoryUpdated cannot be nil")
	}
	r.Signal.Notify(history)
}
