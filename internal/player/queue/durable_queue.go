package queue

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"slices"
	"sync"

	"github.com/RayZ3R0/sonami-gtk/internal/signals"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
)

type durableQueue struct {
	sync.RWMutex
	queue Queue

	source *signals.StatefulSignal[[]sonami.Track]
}

func (q *durableQueue) Append(track sonami.Track) {
	if q.Contains(track.ID()) {
		return
	}

	q.Lock()
	defer q.Unlock()

	q.source.Notify(func(oldValue []sonami.Track) []sonami.Track {
		return append(oldValue, track)
	})
	q.queue.Append(track)
}

func (q *durableQueue) Clear() {
	q.Lock()
	defer q.Unlock()

	q.source.Notify(func(oldValue []sonami.Track) []sonami.Track {
		return []sonami.Track{}
	})
	q.queue.Clear()
}

func (q *durableQueue) Contains(trackID string) bool {
	q.RLock()
	defer q.RUnlock()

	for _, track := range q.source.CurrentValue() {
		if track.ID() == trackID {
			return true
		}
	}
	return false
}

func (q *durableQueue) Entries() *signals.StatefulSignal[[]sonami.Track] {
	return q.queue.Entries()
}

func (q *durableQueue) Get(index int) sonami.Track {
	q.RLock()
	defer q.RUnlock()
	return q.queue.Get(index)
}

func (q *durableQueue) indexOf(trackID string) int {
	for i, track := range q.source.CurrentValue() {
		if track.ID() == trackID {
			return i
		}
	}
	return -1
}

func (q *durableQueue) Insert(track sonami.Track, index int) error {
	if q.Contains(track.ID()) {
		return nil
	}

	q.Lock()
	defer q.Unlock()

	errChan := make(chan error, 1)
	defer close(errChan)

	sourceIndex := q.indexOf(track.ID())
	q.source.Notify(func(oldValue []sonami.Track) []sonami.Track {
		if sourceIndex < 0 || sourceIndex > len(oldValue) {
			errChan <- fmt.Errorf("invalid index, must be between 0 and %d", len(oldValue))
			return oldValue
		}
		errChan <- nil
		return append(oldValue[:sourceIndex], append([]sonami.Track{track}, oldValue[sourceIndex:]...)...)
	})
	if err := <-errChan; err != nil {
		return err
	}

	return q.queue.Insert(track, index)
}

func (q *durableQueue) Peek() sonami.Track {
	q.RLock()
	defer q.RUnlock()
	return q.queue.Peek()
}

func (q *durableQueue) Pop() sonami.Track {
	q.RLock()
	defer q.RUnlock()
	return q.queue.Pop()
}

func (q *durableQueue) Prepend(track sonami.Track) {
	if q.Contains(track.ID()) {
		return
	}

	q.Lock()
	defer q.Unlock()

	q.source.Notify(func(oldValue []sonami.Track) []sonami.Track {
		return append([]sonami.Track{track}, oldValue...)
	})
	q.queue.Prepend(track)
}

func (q *durableQueue) RemoveAt(index int) error {
	q.Lock()
	defer q.Unlock()

	errChan := make(chan error, 1)
	defer close(errChan)

	track := q.queue.Get(index)
	if track == nil {
		return errors.New("track not found")
	}

	sourceIndex := q.indexOf(track.ID())
	q.source.Notify(func(oldValue []sonami.Track) []sonami.Track {
		if sourceIndex < 0 || sourceIndex >= len(oldValue) {
			errChan <- fmt.Errorf("invalid index, must be between 0 and %d", len(oldValue))
			return oldValue
		}
		errChan <- nil
		return append(oldValue[:sourceIndex], oldValue[sourceIndex+1:]...)
	})
	if err := <-errChan; err != nil {
		return err
	}

	return q.queue.RemoveAt(index)
}

func (q *durableQueue) Restore(currentTrackID string) {
	q.Lock()
	defer q.Unlock()

	sourceTracks := slices.Clone(q.source.CurrentValue())
	offset := 0
	for i, track := range sourceTracks {
		if track.ID() == currentTrackID {
			offset = i + 1
			break
		}
	}
	q.queue.Set(sourceTracks[offset:])
}

func (q *durableQueue) Set(tracks []sonami.Track) {
	q.Lock()
	defer q.Unlock()

	q.source.Notify(func(oldValue []sonami.Track) []sonami.Track {
		return slices.Clone(tracks)
	})
	q.queue.Set(tracks)
}

func (q *durableQueue) Shuffle(currentTrackID string) {
	q.Lock()
	defer q.Unlock()

	sourceTracks := slices.Clone(q.source.CurrentValue())
	trackOffset := 0
	for i, track := range sourceTracks {
		if track.ID() == currentTrackID {
			trackOffset = i
			break
		}
	}

	if trackOffset != -1 {
		if trackOffset+1 < len(sourceTracks) {
			sourceTracks = append(sourceTracks[:trackOffset], sourceTracks[trackOffset+1:]...)
		} else {
			sourceTracks = []sonami.Track{}
		}
	}

	rand.Shuffle(len(sourceTracks), func(i, j int) {
		sourceTracks[i], sourceTracks[j] = sourceTracks[j], sourceTracks[i]
	})
	q.queue.Set(sourceTracks)
}

func (q *durableQueue) Skip(n int) ([]sonami.Track, error) {
	q.Lock()
	defer q.Unlock()
	return q.queue.Skip(n)
}

func NewDurableQueue() DurableQueue {
	return &durableQueue{
		queue:  NewQueue(),
		source: signals.NewStatefulSignal([]sonami.Track{}),
	}
}
