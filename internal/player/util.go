package player

import (
	"fmt"

	"codeberg.org/dergs/tonearm/pkg/tonearm"
)

func clearQueues() {
	UserQueue.Clear()
	BaseQueue.Clear()
}

func getNextTrackFromQueue(peek bool) tonearm.Track {
	verb := "pop"
	if peek {
		verb = "peek"
	}
	logger.Debug(fmt.Sprintf("attempting to %s next track from user queue", verb))
	var nextTrack tonearm.Track
	if peek {
		nextTrack = UserQueue.Peek()
	} else {
		nextTrack = UserQueue.Pop()
	}
	if nextTrack != nil {
		logger.Info(fmt.Sprintf("%sed next track from user queue", verb), "track_id", nextTrack.ID())
		return nextTrack
	}

	logger.Debug(fmt.Sprintf("attempting to %s next track from base queue", verb))
	if peek {
		nextTrack = BaseQueue.Peek()
	} else {
		nextTrack = BaseQueue.Pop()
	}
	if nextTrack != nil {
		logger.Info(fmt.Sprintf("%sed next track from base queue", verb), "track_id", nextTrack.ID())
		return nextTrack
	}

	// Check if we are suposed to repeat the base queue
	if RepeatModeChanged.CurrentValue() == RepeatModeQueue {
		logger.Debug("queue repeat mode is enabled, replaying base queue")
		BaseQueue.Restore("")
		if ShuffleStateChanged.CurrentValue() {
			BaseQueue.Shuffle("")
		}
		return BaseQueue.Peek()
	}

	return nil
}

// func resolveTrack(trackId string) (tonearm.Track, error) {
// 	tidal, err := injector.Inject[*tidalapi.TidalAPI]()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return tidal.OpenAPI.V2.Tracks.Track(context.Background(), trackId, "albums.coverArt", "artists")
// }
