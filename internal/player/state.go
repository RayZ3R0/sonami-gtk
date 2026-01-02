package player

import (
	"context"
	"strconv"

	"codeberg.org/dergs/tidalwave/internal/signals"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"github.com/infinytum/injector"
)

func init() {
	OnStateChanged.On(func(state State) bool {
		if OnTrackChanged.current.ID == "" {
			return signals.Continue
		}

		if state.Status == StatusStopped {
			onTrackFinished()
		}
		return signals.Continue
	})
}

func onTrackFinished() {
	defer logger.Debug("done handling track end")
	go nextTrack()
}

func nextTrack() {
	logger.Debug("attempting to pop next track from user queue")
	nextTrack := UserQueue.Pop()
	if nextTrack != nil {
		logger.Info("playing next track from user queue", "track_id", nextTrack.Data.ID)
		playTrack(nextTrack)
		return
	}

	logger.Debug("attempting to pop next track from base queue")
	nextTrack = BaseQueue.Pop()
	if nextTrack != nil {
		logger.Info("playing next track from base queue", "track_id", nextTrack.Data.ID)
		playTrack(nextTrack)
		return
	}

	// Since no other songs are left in the queue, retrieve mix to play from API
	tidal := injector.MustInject[*tidalapi.TidalAPI]()
	trackId, err := strconv.Atoi(OnTrackChanged.current.ID)
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
