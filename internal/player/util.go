package player

import (
	"context"
	"fmt"
	"math/rand/v2"

	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"
	"github.com/infinytum/injector"
)

func clearQueues() {
	UserQueue.Clear()
	BaseQueue.Clear()
}

func getNextTrackFromQueue(peek bool) *openapi.Track {
	verb := "pop"
	if peek {
		verb = "peek"
	}
	logger.Debug(fmt.Sprintf("attempting to %s next track from user queue", verb))
	var nextTrack *openapi.Track
	if peek {
		nextTrack = UserQueue.Peek()
	} else {
		nextTrack = UserQueue.Pop()
	}
	if nextTrack != nil {
		logger.Info(fmt.Sprintf("%sed next track from user queue", verb), "track_id", nextTrack.Data.ID)
		return nextTrack
	}

	logger.Debug(fmt.Sprintf("attempting to %s next track from base queue", verb))
	if peek {
		nextTrack = BaseQueue.Peek()
	} else {
		nextTrack = BaseQueue.Pop()
	}
	if nextTrack != nil {
		logger.Info(fmt.Sprintf("%sed next track from base queue", verb), "track_id", nextTrack.Data.ID)
		return nextTrack
	}

	// Check if we are suposed to repeat the base queue
	if RepeatModeChanged.CurrentValue() == RepeatModeQueue {
		logger.Debug("queue repeat mode is enabled, replaying base queue")
		played := BaseQueue.PastEntries.CurrentValue()
		BaseQueue.Clear()
		for _, track := range played {
			BaseQueue.AddTrack(track, false)
		}
		return BaseQueue.Peek()
	}

	return nil
}

func resolveTrack(trackId string) (*openapi.Track, error) {
	tidal, err := injector.Inject[*tidalapi.TidalAPI]()
	if err != nil {
		return nil, err
	}

	return tidal.OpenAPI.V2.Tracks.Track(context.Background(), trackId, "albums.coverArt", "artists")
}

func prepareTrackList(tracks []openapi.Track, shuffle bool, skipUntilID string) []*openapi.Track {
	if shuffle {
		rand.Shuffle(len(tracks), func(i, j int) {
			tracks[i], tracks[j] = tracks[j], tracks[i]
		})
	} else if skipUntilID != "" {
		for i, track := range tracks {
			if track.Data.ID == skipUntilID {
				tracks = tracks[i:]
				break
			}
		}
	}

	trackPointers := make([]*openapi.Track, len(tracks))
	for i, track := range tracks {
		trackPointers[i] = &track
	}
	return trackPointers
}
