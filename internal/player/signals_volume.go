package player

import (
	"sync"

	"codeberg.org/dergs/tidalwave/internal/signals"
	"codeberg.org/dergs/tidalwave/pkg/mpris"
	"github.com/infinytum/injector"
)

func init() {
	OnVolumeChanged.Signal.On(func(volume float64) bool {
		server := injector.MustInject[*mpris.Server]()
		server.SetVolume(volume)

		return signals.Continue
	})
}

var OnVolumeChanged = volumeChangedSignal{
	signals.NewSignal[func(volume float64) bool](),
	1,
	sync.Mutex{},
}

type volumeChangedSignal struct {
	signals.Signal[func(volume float64) bool]
	current float64
	lock    sync.Mutex
}

func (r *volumeChangedSignal) Notify(callback func(previous float64) float64) {
	r.lock.Lock()
	defer r.lock.Unlock()
	newState := callback(r.current)
	if newState == r.current {
		return
	}
	logger.Debug("volume changed", "volume", newState)
	r.current = newState
	r.Signal.Notify(newState)
}

func (r *volumeChangedSignal) On(handler func(volume float64) bool) *signals.Subscription {
	handler(r.current)
	return r.Signal.On(handler)
}
