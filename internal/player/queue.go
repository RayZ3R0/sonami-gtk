package player

import (
	"log/slog"
	"sync"

	"codeberg.org/dergs/tidalwave/internal/signals"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
)

type Queue struct {
	sync.RWMutex
	logger                 *slog.Logger
	shouldTrackPastEntries bool

	UpcomingEntries *signals.StatefulSignal[[]*openapi.Track]
	PastEntries     *signals.StatefulSignal[[]*openapi.Track]
}

func (q *Queue) AddTrackID(trackId string, addToFront bool) error {
	q.RLock()
	for _, track := range q.UpcomingEntries.CurrentValue() {
		if track.Data.ID == trackId {
			q.RUnlock()
			q.logger.Debug("track is already in queue", "track_id", trackId)
			return nil
		}
	}
	q.RUnlock()

	track, err := resolveTrack(trackId)
	if err != nil {
		return err
	}

	q.AddTrack(track, addToFront)
	return nil
}

func (q *Queue) AddTrack(track *openapi.Track, addToFront bool) {
	q.Lock()
	defer q.Unlock()

	q.UpcomingEntries.Notify(func(oldValue []*openapi.Track) []*openapi.Track {
		if addToFront {
			return append([]*openapi.Track{track}, oldValue...)
		} else {
			return append(oldValue, track)
		}
	})
	q.logger.Info("added track to queue", "track_id", track.Data.ID)
}

func (q *Queue) Clear() {
	q.UpcomingEntries.Notify(func(oldValue []*openapi.Track) []*openapi.Track {
		return []*openapi.Track{}
	})
	q.PastEntries.Notify(func(oldValue []*openapi.Track) []*openapi.Track {
		return []*openapi.Track{}
	})
	q.logger.Info("queue cleared")
}

func (q *Queue) Peek() *openapi.Track {
	q.RLock()
	defer q.RUnlock()

	upcomingTracks := q.UpcomingEntries.CurrentValue()
	if len(upcomingTracks) == 0 {
		return nil
	}
	return upcomingTracks[0]
}

func (q *Queue) Pop() *openapi.Track {
	q.Lock()
	defer q.Unlock()

	upcomingTracks := q.UpcomingEntries.CurrentValue()
	if len(upcomingTracks) == 0 {
		return nil
	}

	nextTrack := upcomingTracks[0]
	q.UpcomingEntries.Notify(func(oldValue []*openapi.Track) []*openapi.Track {
		return oldValue[1:]
	})

	if q.shouldTrackPastEntries {
		q.PastEntries.Notify(func(oldValue []*openapi.Track) []*openapi.Track {
			return append(oldValue, nextTrack)
		})
	}

	return nextTrack
}

func newQueue(logger *slog.Logger, shouldTrackPastEntries bool) *Queue {
	return &Queue{
		RWMutex:                sync.RWMutex{},
		UpcomingEntries:        signals.NewStatefulSignal([]*openapi.Track{}),
		PastEntries:            signals.NewStatefulSignal([]*openapi.Track{}),
		logger:                 logger,
		shouldTrackPastEntries: shouldTrackPastEntries,
	}
}
