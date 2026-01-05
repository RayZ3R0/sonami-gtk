package player

import (
	"context"
	"strconv"

	"codeberg.org/dergs/tidalwave/internal/signals"
	"codeberg.org/dergs/tidalwave/pkg/mpris"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	"github.com/infinytum/injector"
)

func init() {
	OnStateChanged.Signal.On(func(state State) bool {
		mprisClient := injector.MustInject[*mpris.Server]()
		mprisClient.SetPosition(state.Position)
		switch state.Status {
		case StatusBuffering, StatusPaused:
			mprisClient.SetPlaybackStatus(mpris.PlaybackStatusPaused)
		case StatusPlaying:
			mprisClient.SetPlaybackStatus(mpris.PlaybackStatusPlaying)
		case StatusStopped:
			mprisClient.SetPlaybackStatus(mpris.PlaybackStatusStopped)
		}

		return signals.Continue
	})

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

func peekNextTrack() *openapi.Track {
	logger.Debug("attempting to peek next track from user queue")
	nextTrack := UserQueue.Peek()
	if nextTrack != nil {
		logger.Info("peeked next track from user queue", "track_id", nextTrack.Data.ID)
		return nextTrack
	}

	logger.Debug("attempting to peek next track from base queue")
	nextTrack = BaseQueue.Peek()
	if nextTrack != nil {
		logger.Info("peeked next track from base queue", "track_id", nextTrack.Data.ID)
		return nextTrack
	}

	return nil
}

func nextTrack() {
	if OnRepeatModeChanged.current == RepeatModeSingle {
		logger.Debug("single repeat mode is enabled, replaying track")
		Scrub(0)
		return
	}

	logger.Debug("attempting to pop next track from user queue")
	nextTrack := UserQueue.Pop()
	if nextTrack != nil {
		currentHistoryType = HistoryTypeUnmanaged // If the user wants to play their tracks on loop, they should make a playlist.
		logger.Info("playing next track from user queue", "track_id", nextTrack.Data.ID)
		playTrack(nextTrack)
		return
	}

	logger.Debug("attempting to pop next track from base queue")
	nextTrack = BaseQueue.Pop()
	if nextTrack != nil {
		currentHistoryType = HistoryTypeManaged
		logger.Info("playing next track from base queue", "track_id", nextTrack.Data.ID)
		playTrack(nextTrack)
		return
	}

	if OnRepeatModeChanged.current == RepeatModeList {
		logger.Debug("list repeat mode is enabled, replaying history")
		if len(managedHistory.Entries) == 0 {
			logger.Debug("no history to replay, replaying single track")
			Scrub(0)
			return
		}

		firstEntry := *managedHistory.Entries[0]
		playTrackId(firstEntry.TrackID)

		entries := append(managedHistory.Entries[1:], managedHistory.Current)
		managedHistory.Clear()
		managedHistory.Push(&firstEntry)

		for _, entry := range entries {
			BaseQueue.AddTrackID(entry.TrackID, false)
		}

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
	PlayMix(mix.ID)
}
