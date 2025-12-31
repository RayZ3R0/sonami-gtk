package player

import (
	"context"
	"sync"

	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	"github.com/infinytum/injector"
)

var userQueue []*openapi.Track
var userQueueMutex sync.RWMutex

var addToQueueMutex sync.Mutex

func AddToQueue(trackId string, addToFront bool) error {
	addToQueueMutex.Lock()
	defer addToQueueMutex.Unlock()

	// If we added a song to the queue and nothing is playing, the user likely wants to start playing the queue
	if OnStateChanged.current.Status == StatusStopped {
		logger.Info("no track is currently playing, immediately playing track", "track_id", trackId)
		Play(trackId)
		return nil
	}

	userQueueMutex.RLock()
	for _, track := range userQueue {
		if track.Data.ID == trackId {
			userQueueMutex.RUnlock()
			logger.Debug("track is already in user-queue", "track_id", trackId)
			return nil
		}
	}
	userQueueMutex.RUnlock()

	tidal, err := injector.Inject[*tidalapi.TidalAPI]()
	if err != nil {
		return err
	}

	openTrack, err := tidal.OpenAPI.V2.Tracks.Track(context.Background(), trackId, "albums.coverArt", "artists")
	if err != nil {
		return err
	}

	modifyUserQueue(func(t []*openapi.Track) []*openapi.Track {
		if addToFront {
			return append([]*openapi.Track{openTrack}, t...)
		} else {
			return append(t, openTrack)
		}
	})
	logger.Info("added track to user-queue", "track_id", trackId)

	return nil
}

func PopFromQueue() *openapi.Track {
	userQueueMutex.RLock()
	if len(userQueue) == 0 {
		userQueueMutex.RUnlock()
		logger.Debug("pop from user-queue failed, queue is empty")
		return nil
	}
	nextTrack := userQueue[0]
	userQueueMutex.RUnlock()

	modifyUserQueue(func(t []*openapi.Track) []*openapi.Track {
		return t[1:]
	})
	logger.Debug("popped track from user-queue", "track_id", nextTrack.Data.ID)

	return nextTrack
}

func modifyUserQueue(callback func([]*openapi.Track) []*openapi.Track) {
	userQueueMutex.Lock()
	userQueue = callback(userQueue)
	userQueueMutex.Unlock()
	OnUserQueueChanged.Notify(func() []*openapi.Track {
		return userQueue
	})
}
