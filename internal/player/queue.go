package player

import (
	"log/slog"
	"sync"

	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"
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
	_, nextTrack, _ := q.Skip(1)

	return nextTrack
}

func (q *Queue) Remove(index int) {
	q.Lock()
	defer q.Unlock()

	upcomingTracks := q.UpcomingEntries.CurrentValue()
	if index < 0 || index >= len(upcomingTracks) {
		return
	}

	q.UpcomingEntries.Notify(func(oldValue []*openapi.Track) []*openapi.Track {
		return append(oldValue[:index], oldValue[index+1:]...)
	})
}

func (q *Queue) SetTracks(tracks []*openapi.Track) {
	q.Lock()
	defer q.Unlock()

	q.UpcomingEntries.Notify(func(oldValue []*openapi.Track) []*openapi.Track {
		return tracks
	})
}

func (q *Queue) Skip(count int) (ok bool, current *openapi.Track, skipped []*openapi.Track) {
	q.Lock()
	defer q.Unlock()

	upcomingTracks := q.UpcomingEntries.CurrentValue()

	if count <= 0 || count > len(upcomingTracks) {
		return false, nil, nil
	}

	current = upcomingTracks[count-1]

	q.UpcomingEntries.Notify(func(oldValue []*openapi.Track) []*openapi.Track {
		return oldValue[count:]
	})

	newHistoryElements := upcomingTracks[:count]

	if q.shouldTrackPastEntries {
		q.PastEntries.Notify(func(oldValue []*openapi.Track) []*openapi.Track {
			return append(oldValue, newHistoryElements...)
		})
	}

	skipped = newHistoryElements[:len(newHistoryElements)-1]
	ok = true
	return
}

func SkipThoughQueue(queue *Queue, to int) {
	go func() {
		setLoadingState()
		if ok, toPlay, historyElements := BaseQueue.Skip(to + 1); ok {
			playTrack(toPlay)

			for _, track := range historyElements {
				history.Push(&HistoryEntry{track.Data.ID})
			}

			history.Push(&HistoryEntry{toPlay.Data.ID})
		}
	}()
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
