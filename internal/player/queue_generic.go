package player

import (
	"context"
	"sync"

	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	"github.com/infinytum/injector"
)

var UserQueue = NewQueue(OnUserQueueChanged.Notify)
var BaseQueue = NewQueue(OnBaseQueueChanged.Notify)

type Queue struct {
	sync.RWMutex

	addTrackMutex sync.Mutex
	popTrackMutex sync.Mutex

	entries []*openapi.Track
	notify  func(func() []*openapi.Track)
}

func (q *Queue) AddTrackID(trackId string, addToFront bool) error {
	q.RLock()
	for _, track := range q.entries {
		if track.Data.ID == trackId {
			q.RUnlock()
			logger.Debug("track is already in queue", "track_id", trackId)
			return nil
		}
	}
	q.RUnlock()

	tidal, err := injector.Inject[*tidalapi.TidalAPI]()
	if err != nil {
		return err
	}

	openTrack, err := tidal.OpenAPI.V2.Tracks.Track(context.Background(), trackId, "albums.coverArt", "artists")
	if err != nil {
		return err
	}

	q.AddTrack(openTrack, addToFront)
	return nil
}

func (q *Queue) AddTrack(track *openapi.Track, addToFront bool) {
	q.addTrackMutex.Lock()
	defer q.addTrackMutex.Unlock()

	// If we added a song to the queue and nothing is playing, the user likely wants to start playing the queue
	if OnStateChanged.current.Status == StatusStopped {
		logger.Info("no track is currently playing, immediately playing track", "track_id", track.Data.ID)
		playTrack(track)
		return
	}

	q.ModifyQueue(func(t []*openapi.Track) []*openapi.Track {
		if addToFront {
			return append([]*openapi.Track{track}, t...)
		} else {
			return append(t, track)
		}
	})
	logger.Info("added track to queue", "track_id", track.Data.ID)
}

func (q *Queue) Clear() {
	q.ModifyQueue(func(t []*openapi.Track) []*openapi.Track {
		return []*openapi.Track{}
	})
	logger.Info("queue cleared")
}

func (q *Queue) ModifyQueue(f func(t []*openapi.Track) []*openapi.Track) {
	q.Lock()
	q.entries = f(q.entries)
	q.Unlock()
	q.notify(func() []*openapi.Track {
		return q.entries
	})
}

func (q *Queue) Pop() *openapi.Track {
	q.popTrackMutex.Lock()
	defer q.popTrackMutex.Unlock()

	q.RLock()
	if len(q.entries) == 0 {
		q.RUnlock()
		logger.Debug("pop from queue failed, queue is empty")
		return nil
	}
	nextTrack := q.entries[0]
	q.RUnlock()

	q.ModifyQueue(func(t []*openapi.Track) []*openapi.Track {
		return t[1:]
	})
	logger.Debug("popped track from queue", "track_id", nextTrack.Data.ID)

	return nextTrack
}

func NewQueue(notify func(func() []*openapi.Track)) *Queue {
	return &Queue{
		entries: make([]*openapi.Track, 0),
		notify:  notify,
	}
}
