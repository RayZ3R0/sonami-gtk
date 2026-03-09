package queue

import (
	"fmt"
	"slices"
	"sync"

	"github.com/RayZ3R0/sonami-gtk/internal/signals"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
)

type queue struct {
	sync.RWMutex

	entries *signals.StatefulSignal[[]sonami.Track]
}

func (q *queue) Append(track sonami.Track) {
	q.Lock()
	defer q.Unlock()

	q.entries.Notify(func(oldValue []sonami.Track) []sonami.Track {
		return append(oldValue, track)
	})
}

func (q *queue) Clear() {
	q.Lock()
	defer q.Unlock()

	q.entries.Notify(func(oldValue []sonami.Track) []sonami.Track {
		return []sonami.Track{}
	})
}

func (q *queue) Contains(trackID string) bool {
	q.RLock()
	defer q.RUnlock()

	for _, track := range q.entries.CurrentValue() {
		if track.ID() == trackID {
			return true
		}
	}
	return false
}

func (q *queue) Entries() *signals.StatefulSignal[[]sonami.Track] {
	return q.entries
}

func (q *queue) Get(index int) sonami.Track {
	q.RLock()
	defer q.RUnlock()

	if index < 0 || index >= len(q.entries.CurrentValue()) {
		return nil
	}
	return q.entries.CurrentValue()[index]
}

func (q *queue) Insert(track sonami.Track, index int) error {
	q.Lock()
	defer q.Unlock()

	errChan := make(chan error, 1)
	defer close(errChan)

	q.entries.Notify(func(oldValue []sonami.Track) []sonami.Track {
		if index < 0 || index > len(oldValue) {
			errChan <- fmt.Errorf("invalid index, must be between 0 and %d", len(oldValue))
			return oldValue
		}
		errChan <- nil
		return append(oldValue[:index], append([]sonami.Track{track}, oldValue[index:]...)...)
	})
	return <-errChan
}

func (q *queue) Peek() sonami.Track {
	q.RLock()
	defer q.RUnlock()

	entries := q.entries.CurrentValue()
	if len(entries) == 0 {
		return nil
	}
	return entries[0]
}

func (q *queue) Pop() sonami.Track {
	q.Lock()
	defer q.Unlock()

	trackChan := make(chan sonami.Track, 1)
	defer close(trackChan)

	q.entries.Notify(func(oldValue []sonami.Track) []sonami.Track {
		if len(oldValue) == 0 {
			trackChan <- nil
			return oldValue
		}

		track := oldValue[0]
		trackChan <- track
		return oldValue[1:]
	})

	return <-trackChan
}

func (q *queue) Prepend(track sonami.Track) {
	q.Lock()
	defer q.Unlock()

	q.entries.Notify(func(oldValue []sonami.Track) []sonami.Track {
		return append([]sonami.Track{track}, oldValue...)
	})
}

func (q *queue) RemoveAt(index int) error {
	q.Lock()
	defer q.Unlock()

	errChan := make(chan error, 1)
	defer close(errChan)

	q.entries.Notify(func(oldValue []sonami.Track) []sonami.Track {
		if index < 0 || index >= len(oldValue) {
			errChan <- fmt.Errorf("invalid index, must be between 0 and %d", len(oldValue))
			return oldValue
		}
		errChan <- nil
		return append(oldValue[:index], oldValue[index+1:]...)
	})
	return <-errChan
}

func (q *queue) Set(tracks []sonami.Track) {
	q.Lock()
	defer q.Unlock()

	q.entries.Notify(func(oldValue []sonami.Track) []sonami.Track {
		return slices.Clone(tracks)
	})
}

func (q *queue) Skip(n int) ([]sonami.Track, error) {
	q.Lock()
	defer q.Unlock()

	if n < 0 {
		return nil, fmt.Errorf("invalid number of tracks to skip")
	}

	errChan := make(chan error, 1)
	defer close(errChan)

	skippedTracksChan := make(chan []sonami.Track, 1)
	defer close(skippedTracksChan)

	q.entries.Notify(func(oldValue []sonami.Track) []sonami.Track {
		if len(oldValue) < n {
			errChan <- fmt.Errorf("not enough tracks in queue")
			return oldValue
		}
		errChan <- nil

		skippedTracksChan <- oldValue[:n]
		return oldValue[n:]
	})

	if err := <-errChan; err != nil {
		return nil, err
	}

	return <-skippedTracksChan, nil
}

func NewQueue() Queue {
	return &queue{
		entries: signals.NewStatefulSignal([]sonami.Track{}),
	}
}
