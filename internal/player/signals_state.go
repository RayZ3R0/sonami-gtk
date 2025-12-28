package player

import (
	"sync"

	"codeberg.org/dergs/tidalwave/internal/signals"
)

var OnStateChanged = stateChangedSignal{
	signals.NewSignal[func(state State) bool](),
	State{
		Status: StatusStopped,
	},
	sync.Mutex{},
}

type stateChangedSignal struct {
	signals.Signal[func(state State) bool]
	current State
	lock    sync.Mutex
}

func (r *stateChangedSignal) Notify(callback func(state *State)) {
	r.lock.Lock()
	defer r.lock.Unlock()
	newState := r.current
	callback(&newState)
	if newState.Equals(r.current) {
		return
	}
	logger.Debug("state changed", "duration", newState.Duration, "position", newState.Position, "status", newState.Status)
	r.current = newState
	r.Signal.Notify(newState)
}

func (r *stateChangedSignal) On(handler func(state State) bool) *signals.Subscription {
	handler(r.current)
	return r.Signal.On(handler)
}

type State struct {
	Duration int
	Position int
	Status   Status
}

func (s State) Equals(other State) bool {
	return s.Status == other.Status && s.Duration == other.Duration && s.Position == other.Position
}

type Status string

const (
	StatusBuffering Status = "buffering"
	StatusPlaying   Status = "playing"
	StatusPaused    Status = "paused"
	StatusStopped   Status = "stopped"
)
