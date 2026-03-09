package player

import (
	"fmt"

	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
	"github.com/infinytum/injector"
)

func clearQueues() {
	UserQueue.Clear()
	BaseQueue.Clear()
}

func getNextTrackFromQueue(peek bool) sonami.Track {
	verb := "pop"
	if peek {
		verb = "peek"
	}
	logger.Debug(fmt.Sprintf("attempting to %s next track from user queue", verb))
	var nextTrack sonami.Track
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

func resolveTrack(trackId string) (sonami.Track, error) {
	service, err := injector.Inject[sonami.Service]()
	if err != nil {
		return nil, err
	}

	return service.GetTrack(trackId)
}
