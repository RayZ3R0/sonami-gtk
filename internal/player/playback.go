package player

import (
	"context"
	"strconv"

	"codeberg.org/dergs/tonearm/internal/settings"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	v1 "codeberg.org/dergs/tonearm/pkg/tidalapi/models/v1"
	tracksv1 "codeberg.org/dergs/tonearm/pkg/tidalapi/v1/tracks"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
	"github.com/go-gst/go-gst/gst"
	"github.com/infinytum/injector"
)

func playTrack(track tonearm.Track) error {
	tidal, err := injector.Inject[*tidalapi.TidalAPI]()
	if err != nil {
		return err
	}

	TrackChanged.Notify(func(oldState tonearm.Track) tonearm.Track {
		return track
	})

	PlaybackStateChanged.Notify(func(oldValue *PlaybackState) *PlaybackState {
		newState := *oldValue
		newState.Duration = track.Duration()
		return &newState
	})

	// if !slices.Contains(track.Data.Attributes.Availability, openapi.TrackAvailabilityStream) {
	// 	notifications.OnToast.Notify(gettext.Get("Track not available for streaming, skipping to next track"))
	// 	Next()
	// 	return errors.New("track not available for streaming")
	// }

	if currentlyEnqueuedTrack == nil || strconv.Itoa(currentlyEnqueuedTrack.TrackID) != track.ID() {
		logger.Debug("fetching playback info for track", "track_id", track.ID())
		playbackInfo, err := tidal.V1.Tracks.PlaybackInfo(
			context.Background(),
			track.ID(),
			tracksv1.PlaybackInfoOptions{
				AudioQuality: settings.Player().GetAudioQuality(),
			},
		)
		if err != nil {
			logger.Error("unable to fetch playback info for track", "error", err)
			return err
		}
		return play(playbackInfo)
	}
	logger.Debug("gapless playback detected, not enqueueing track again")
	resetLoadingState()
	return nil
}

func play(playbackInfo *v1.PlaybackInfo) error {
	// Inform the UI about the track quality we got from TIDAL.
	PlaybackQualityChanged.Notify(func(oldValue v1.AudioQuality) v1.AudioQuality {
		return playbackInfo.AudioQuality
	})

	// Free up resources taken up by previous stream
	playbin.SetState(gst.StateNull)
	playbin.SetArg("uri", "")

	PlaybackStateChanged.Notify(func(oldValue *PlaybackState) *PlaybackState {
		newState := *oldValue
		newState.Loading = true
		return &newState
	})

	if err := enqueue(playbackInfo); err != nil {
		return err
	}

	PlaybackStateChanged.Notify(func(oldValue *PlaybackState) *PlaybackState {
		newState := *oldValue
		newState.Status = PlaybackStatusPlaying
		newState.Loading = false
		return &newState
	})

	playbin.SetState(gst.StatePlaying)
	return nil
}
