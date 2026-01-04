package player

import (
	"sync"

	"codeberg.org/dergs/tidalwave/internal/signals"
	v1 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v1"
)

var OnPlaybackQualityChanged = playbackQualityChangedSignal{
	signals.NewSignal[func(quality v1.AudioQuality) bool](),
	v1.AudioQualityHighResLossless,
	sync.Mutex{},
}

type playbackQualityChangedSignal struct {
	signals.Signal[func(quality v1.AudioQuality) bool]
	current v1.AudioQuality
	lock    sync.Mutex
}

func (r *playbackQualityChangedSignal) Notify(callback func() v1.AudioQuality) {
	r.lock.Lock()
	defer r.lock.Unlock()
	newState := callback()
	if newState == r.current {
		return
	}
	logger.Info("playback quality changed", "quality", newState)
	r.current = newState
	r.Signal.Notify(newState)
}

func (r *playbackQualityChangedSignal) On(handler func(quality v1.AudioQuality) bool) *signals.Subscription {
	handler(r.current)
	return r.Signal.On(handler)
}
