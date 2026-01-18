package player

import (
	"time"

	"github.com/go-gst/go-gst/gst"
)

var (
	BaseQueue = newQueue(logger, true)
	UserQueue = newQueue(logger, false)
)

func playNextTrack() {
	if RepeatModeChanged.CurrentValue() == RepeatModeTrack {
		logger.Debug("single repeat mode is enabled, replaying track")
		SeekToPosition(0)
		startUpdateRunner()
		return
	}

	setLoadingState()
	nextTrack := getNextTrackFromQueue(false)
	if nextTrack != nil {
		logger.Info("playing next track", "track_id", nextTrack.Data.ID)
		playTrack(nextTrack)
		history.Push(&HistoryEntry{
			TrackID: nextTrack.Data.ID,
		})

		return
	}

	// Since no other songs are left in the queue, retrieve mix to play from API
	logger.Info("starting track radio", "track_id", TrackChanged.CurrentValue().ID)
	PlayTrackRadio(TrackChanged.CurrentValue().ID, true)
}

func playPreviousTrack() {
	ok, position := playbin.QueryPosition(gst.FormatTime)
	if ok && time.Duration(position) > 5*time.Second {
		logger.Debug("above the 5 second mark, replaying song", "action", "previous")
		SeekToPosition(0)
		return
	}

	if len(history.Entries.CurrentValue()) < 1 {
		logger.Debug("no history entries, replaying song", "action", "previous")
		SeekToPosition(0)
		return
	}

	setLoadingState()

	entry := history.Pop()
	if entry != nil {
		// Re-Queue current song to front of user-queue
		UserQueue.AddTrackID(TrackChanged.CurrentValue().ID, true)

		// Switch to previous track without clearing base queue
		track, err := resolveTrack(entry.TrackID)
		if err != nil {
			logger.Error("failed to resolve track", "trackID", entry.TrackID, "error", err)
			return
		}
		logger.Debug("playing previous track", "track_id", track.Data.ID, "action", "previous")
		playTrack(track)
	}
}
