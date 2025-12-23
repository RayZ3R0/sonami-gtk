package player

import (
	"sync"

	"codeberg.org/dergs/tidalwave/internal/ui/signals"
	v1 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v1"
)

var OnState = stateSignal{
	signals.NewSignal[func(state State) bool](),
	State{
		Status: StatusStopped,
		Track:  nil,
	},
	sync.Mutex{},
}

type stateSignal struct {
	signals.Signal[func(state State) bool]
	current State
	lock    sync.Mutex
}

func (r *stateSignal) Notify(callback func(state *State)) {
	r.lock.Lock()
	defer r.lock.Unlock()
	newState := r.current
	callback(&newState)
	if newState.Equals(r.current) {
		return
	}
	r.current = newState
	r.Signal.Notify(newState)
}

func (r *stateSignal) On(handler func(state State) bool) *signals.Subscription {
	handler(r.current)
	return r.Signal.On(handler)
}

type State struct {
	Duration int
	Position int
	Status   Status
	Track    *v1.Track
	Volume   float64
}

func (s State) Equals(other State) bool {
	return s.Status == other.Status && s.Track == other.Track && s.Duration == other.Duration && s.Position == other.Position && s.Volume == other.Volume
}

type Status string

const (
	StatusBuffering Status = "buffering"
	StatusPlaying   Status = "playing"
	StatusPaused    Status = "paused"
	StatusStopped   Status = "stopped"
)
