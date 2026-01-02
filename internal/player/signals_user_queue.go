package player

import (
	"sync"

	"codeberg.org/dergs/tidalwave/internal/signals"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
)

var OnBaseQueueChanged = userQueueChangedSignal{
	signals.NewSignal[func(tracks []*openapi.Track) bool](),
	[]*openapi.Track{},
	sync.Mutex{},
}

var OnUserQueueChanged = userQueueChangedSignal{
	signals.NewSignal[func(tracks []*openapi.Track) bool](),
	[]*openapi.Track{},
	sync.Mutex{},
}

type userQueueChangedSignal struct {
	signals.Signal[func(tracks []*openapi.Track) bool]
	current []*openapi.Track
	lock    sync.Mutex
}

func (r *userQueueChangedSignal) Notify(callback func() []*openapi.Track) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.current = callback()
	r.Signal.Notify(r.current)
}

func (r *userQueueChangedSignal) On(handler func(tracks []*openapi.Track) bool) *signals.Subscription {
	handler(r.current)
	return r.Signal.On(handler)
}
