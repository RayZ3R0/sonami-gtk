package player

import (
	"context"
	"strconv"
	"time"

	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"github.com/go-gst/go-gst/gst"
	"github.com/infinytum/injector"
)

var (
	BaseQueue = newQueue(logger, true)
	UserQueue = newQueue(logger, false)
)

func playNextTrack() {
	if RepeatModeChanged.CurrentValue() == RepeatModeTrack {
		logger.Debug("single repeat mode is enabled, replaying track")
		SeekToPosition(0)
		return
	}

	nextTrack := getNextTrackFromQueue(false)
	if nextTrack != nil {
		logger.Info("playing next track", "track_id", nextTrack.Data.ID)
		playTrack(nextTrack)
		return
	}

	// Since no other songs are left in the queue, retrieve mix to play from API
	tidal := injector.MustInject[*tidalapi.TidalAPI]()
	trackId, err := strconv.Atoi(TrackChanged.CurrentValue().ID)
	if err != nil {
		logger.Error("failed to parse album id", "error", err)
		return
	}

	mix, err := tidal.V1.Tracks.Mix(context.Background(), trackId)
	if err != nil {
		logger.Error("failed to retrieve mix", "error", err)
		return
	}
	logger.Info("starting track radio", "mix_id", mix.ID)
	PlayPlaylist(mix.ID, false, "")
}

func playPreviousTrack() {
	ok, position := playbin.QueryPosition(gst.FormatTime)
	if ok && time.Duration(position) > 5*time.Second {
		SeekToPosition(0)
		return
	}

	if len(history.Entries.CurrentValue()) < 1 {
		SeekToPosition(0)
		return
	}

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
		playTrack(track)
	}
}
