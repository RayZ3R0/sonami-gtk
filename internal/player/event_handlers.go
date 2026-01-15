package player

import (
	"context"
	"strconv"
	"time"

	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	v1 "codeberg.org/dergs/tonearm/pkg/tidalapi/models/v1"
	tracksv1 "codeberg.org/dergs/tonearm/pkg/tidalapi/v1/tracks"
	"github.com/go-gst/go-gst/gst"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/glib"
)

const UpdateInterval = 250 * time.Millisecond

var updateRunnerSourceHandle uint

func onAboutToFinish(_ *gst.Element) {
	if RepeatModeChanged.CurrentValue() == RepeatModeTrack {
		return
	}

	nextTrack := getNextTrackFromQueue(true)
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
		TrackChanged.OnLazy(func(t *Track) bool {
			logger.Debug("triggered one-shot handler to propagate gapless playback quality")
			PlaybackQualityChanged.Notify(func(oldValue v1.AudioQuality) v1.AudioQuality {
				return playbackInfo.AudioQuality
			})
			return signals.Unsubscribe
		})
	}
}

func onBusMessage(msg *gst.Message) bool {
	switch msg.Type() {
	case gst.MessageError:
		err := msg.ParseError()
		logger.Error("playback failed", "code", err.Code(), "message", err.Message(), "error", err.Error(), "debug", err.DebugString())
	case gst.MessageStreamStart:
		startUpdateRunner()
		playbin.SetProperty("volume", volumeBeforeEnqueue)
		// A hack to trigger the correct track updates with gapless playback
		if TrackChanged.CurrentValue().ID != strconv.Itoa(currentlyEnqueuedTrackID) {
			playNextTrack()
		}
	case gst.MessageEOS:
		stopUpdateRunner()
		PlaybackStateChanged.Notify(func(oldValue *PlaybackState) *PlaybackState {
			newState := *oldValue
			newState.Status = PlaybackStatusStopped
			return &newState
		})
		go playNextTrack()
	case gst.MessageStateChanged:
		_, newState := msg.ParseStateChanged()
		switch newState {
		case gst.StatePlaying:
			PlaybackStateChanged.Notify(func(oldValue *PlaybackState) *PlaybackState {
				newState := *oldValue
				newState.Status = PlaybackStatusPlaying
				return &newState
			})
		}
	}
	return true
}

func onVolumeChange() {
	if volume, err := playbin.GetProperty("volume"); err == nil {
		return
	} else if volume.(float64) == VolumeChanged.CurrentValue() {
		return
	} else {
		VolumeChanged.Notify(func(oldValue float64) float64 {
			return volume.(float64)
		})
	}
}

var onUpdateTick = glib.SourceFunc(func(uintptr) bool {
	PlaybackStateChanged.Notify(func(oldValue *PlaybackState) *PlaybackState {
		newState := *oldValue
		if ok, duration := playbin.QueryDuration(gst.FormatTime); ok {
			newState.Duration = time.Duration(duration)
		}

		if ok, position := playbin.QueryPosition(gst.FormatTime); ok {
			newState.Position = time.Duration(position)
			newState.IsSeeking = false
		}
		return &newState
	})
	return glib.SOURCE_CONTINUE
})

func startUpdateRunner() {
	if updateRunnerSourceHandle != 0 {
		return
	}
	updateRunnerSourceHandle = glib.TimeoutAdd(uint(UpdateInterval.Milliseconds()), &onUpdateTick, 0)
	logger.Debug("started playbin update runner", "source_handle", updateRunnerSourceHandle)
}

func stopUpdateRunner() {
	if updateRunnerSourceHandle == 0 {
		return
	}

	glib.SourceRemove(updateRunnerSourceHandle)
	logger.Debug("stopped playbin update runner", "source_handle", updateRunnerSourceHandle)
	updateRunnerSourceHandle = 0
}
