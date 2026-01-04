package player

import (
	"context"

	"codeberg.org/dergs/tidalwave/internal/signals"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	v1 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v1"
	tracksv1 "codeberg.org/dergs/tidalwave/pkg/tidalapi/v1/tracks"
	"github.com/go-gst/go-gst/gst"
	"github.com/infinytum/injector"
)

func onAboutToFinish(_ *gst.Element) {
	nextTrack := peekNextTrack()
	if nextTrack != nil {
		playbackInfo, err := injector.MustInject[*tidalapi.TidalAPI]().V1.Tracks.PlaybackInfo(context.Background(), nextTrack.Data.ID, tracksv1.PlaybackInfoOptions{})
		if err != nil {
			logger.Error("failed to get playback info", "error", err)
			return
		}
		if err := enqueue(playbackInfo); err != nil {
			logger.Error("enqueueing for gapless playback", "error", err)
			return
		}
		logger.Info("enqueued next song for gapless playback", "track_id", nextTrack.Data.ID)

		// One-Shot Handler to update the track quality
		OnTrackChanged.On(func(trackInfo TrackInformation) bool {
			logger.Debug("triggered one-shot handler to propagate gapless playback quality")
			OnPlaybackQualityChanged.Notify(func() v1.AudioQuality {
				return playbackInfo.AudioQuality
			})
			return signals.Unsubscribe
		})
	}
}
