package player

import (
	"context"
	"fmt"
	"time"

	"codeberg.org/dergs/tonearm/internal/g"
	"codeberg.org/dergs/tonearm/internal/settings"
	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	v1 "codeberg.org/dergs/tonearm/pkg/tidalapi/models/v1"
	tracksv1 "codeberg.org/dergs/tonearm/pkg/tidalapi/v1/tracks"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
	"github.com/go-gst/go-gst/gst"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/glib"
)

const UpdateInterval = 250 * time.Millisecond

var updateRunnerSourceHandle uint
var didQueueGaplessPlayback bool
var disableGaplessPlayback bool

func onAboutToFinish(_ *gst.Element) {
	if RepeatModeChanged.CurrentValue() == RepeatModeTrack {
		return
	}

	if disableGaplessPlayback {
		return
	}

	nextTrack := getNextTrackFromQueue(true)
	if nextTrack != nil {
		playbackInfo, err := injector.MustInject[*tidalapi.TidalAPI]().V1.Tracks.PlaybackInfo(
			context.Background(),
			nextTrack.ID(),
			tracksv1.PlaybackInfoOptions{
				AudioQuality: settings.Player().GetAudioQuality(),
			},
		)
		if err != nil {
			logger.Error("failed to get playback info", "error", err)
			return
		}
		if err := enqueue(playbackInfo); err != nil {
			logger.Error("enqueueing for gapless playback", "error", err)
			return
		}
		logger.Info("enqueued next song for gapless playback", "track_id", nextTrack.ID())
		didQueueGaplessPlayback = true

		// One-Shot Handler to update the track quality
		TrackChanged.OnLazy(func(t tonearm.Track) bool {
			logger.Debug("triggered one-shot handler to propagate gapless playback quality")
			PlaybackQualityChanged.Notify(func(oldValue v1.AudioQuality) v1.AudioQuality {
				return playbackInfo.AudioQuality
			})
			return signals.Unsubscribe
		})
	}
}

type codec string

const (
	flac codec = "Free Lossless Audio Codec (FLAC)"
	aac  codec = "MPEG-4 AAC"
)

func (c codec) String() string {
	switch c {
	case flac:
		return "FLAC"
	case aac:
		return "AAC"
	default:
		return "unknown codec"
	}
}

func onBusMessage(msg *gst.Message) bool {
	switch msg.Type() {
	case gst.MessageError:
		err := msg.ParseError()
		logger.Error("playback failed", "code", err.Code(), "message", err.Message(), "error", err.Error(), "debug", err.DebugString())
	case gst.MessageTag:
		// The logs in this branch are all debug because of the nature of MessageTag.
		// If they weren't Debug, the console would get spammed with error messages,
		// since it is expected to have missing data. Between streams

		// CODEC
		tagList := msg.ParseTags()
		if tagList == nil {
			logger.Debug("Error while getting codec")
			return true
		}

		codecIdentifier, ok := tagList.GetString(gst.TagAudioCodec)
		if !ok {
			logger.Debug("Error while getting codec")
			return true
		}
		codec := codec(codecIdentifier)

		// AUDIO CAPS
		format, rate, err := getAudioCaps()
		if err != nil {
			logger.Debug("Error while getting audio stream quality", "error", err)
			return true
		}

		/// Format
		var readableBitDepth string
		switch format {
		case "S16LE":
			readableBitDepth = "16-bit"
		case "S24LE", "S24_32LE":
			readableBitDepth = "24-bit"
		case "S32LE":
			readableBitDepth = "32-bit"
		}

		// Compilation
		var res string
		switch codec {
		case aac:
			// BITRATE (for AAC)
			rate, ok := tagList.GetUint32(gst.TagMaximumBitrate)
			if !ok {
				logger.Debug("Error while getting bitrate")
				return true
			}

			bitrate := g.TruncateFloat(float64(rate)/1000, 1)
			res = fmt.Sprintf("%s %s kbps AAC", readableBitDepth, bitrate)
		case flac:
			/// Sample rate (FLAC)
			sampleRate := g.TruncateFloat(float64(rate)/1000, 1)

			res = fmt.Sprintf("%s %skHz FLAC", readableBitDepth, sampleRate)
		}

		AudioStreamQuality.SetValue(res)
	case gst.MessageStreamStart:
		startUpdateRunner()
		playbin.Set("volume", settings.Player().GetVolume())
		// A hack to trigger the correct track updates with gapless playback
		if didQueueGaplessPlayback {
			stateBeforeLoading = gst.StatePlaying
			go playNextTrack()
		}
	case gst.MessageEOS:
		playbin.SetState(gst.StatePaused)
		stopUpdateRunner()
		PlaybackStateChanged.Notify(func(oldValue *PlaybackState) *PlaybackState {
			newState := *oldValue
			newState.Status = PlaybackStatusStopped
			return &newState
		})
		go playNextTrack()
	case gst.MessageBuffering:
		percent := msg.ParseBuffering()
		if percent == 100 {
			playbin.SetState(gst.StatePlaying)
		} else {
			playbin.SetState(gst.StatePaused)
		}
	case gst.MessageStateChanged:
		_, newState := msg.ParseStateChanged()
		switch newState {
		case gst.StatePlaying:
			startUpdateRunner()
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
	if volume, err := playbin.GetProperty("volume"); err != nil {
		return
	} else if volume.(float64) == VolumeChanged.CurrentValue() {
		return
	} else {
		settings.Player().SetVolume(volume.(float64))
		VolumeChanged.Notify(func(oldValue float64) float64 {
			return volume.(float64)
		})
	}
}

var onUpdateTick = glib.SourceFunc(func(uintptr) bool {
	PlaybackStateChanged.Notify(func(oldValue *PlaybackState) *PlaybackState {
		newState := *oldValue
		if ok, duration := playbin.QueryDuration(gst.FormatTime); ok && !didQueueGaplessPlayback {
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
