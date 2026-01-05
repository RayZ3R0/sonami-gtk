package player

import (
	"sync"

	"codeberg.org/dergs/tidalwave/internal/signals"
)

type RepeatMode int

const (
	RepeatModeNone RepeatMode = iota
	RepeatModeList
	RepeatModeSingle
)

var OnRepeatModeChanged = repeatModeSignal{
	signals.NewSignal[func(state RepeatMode) bool](),
	RepeatModeNone,
	sync.Mutex{},
}

type repeatModeSignal struct {
	signals.Signal[func(state RepeatMode) bool]
	current RepeatMode
	lock    sync.Mutex
}

func (r *repeatModeSignal) Notify(callback func(oldMode *RepeatMode)) {
	r.lock.Lock()
	defer r.lock.Unlock()
	newMode := r.current
	callback(&newMode)
	if newMode == r.current {
		return
	}
	logger.Debug("repeat mode changed", "new mode", newMode)
	r.current = newMode
	r.Signal.Notify(newMode)
}

func (r *repeatModeSignal) On(handler func(state RepeatMode) bool) *signals.Subscription {
	handler(r.current)
	return r.Signal.On(handler)
}
