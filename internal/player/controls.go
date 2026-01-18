package player

import (
	"fmt"
	"math"
	"sync"
	"time"

	"codeberg.org/dergs/tonearm/internal/signals"
	"github.com/go-gst/go-gst/gst"
)

type ControllableState struct {
	HasTrack    bool
	PlayerReady bool
}

func (cs *ControllableState) CanControl() bool {
	return cs.HasTrack && cs.PlayerReady
}

func init() {
	PlaybackStateChanged.On(func(ps *PlaybackState) bool {
		playerReady := ps.Status != PlaybackStatusBuffering &&
			ps.Status != PlaybackStatusLoadingTrack

		if playerReady != ControllableStateChanged.CurrentValue().PlayerReady {
			ControllableStateChanged.Notify(func(oldValue ControllableState) ControllableState {
				oldValue.PlayerReady = playerReady
				return oldValue
			})
		}

		return signals.Continue
	})

	TrackChanged.On(func(newTrack *Track) bool {
		hasTrack := newTrack != nil
		if hasTrack != ControllableStateChanged.CurrentValue().HasTrack {
			ControllableStateChanged.Notify(func(oldValue ControllableState) ControllableState {
				oldValue.HasTrack = hasTrack
				return oldValue
			})
		}

		return signals.Continue
	})
}

func CycleRepeatMode() {
	switch RepeatModeChanged.CurrentValue() {
	case RepeatModeNone:
		SetRepeatMode(RepeatModeQueue)
	case RepeatModeQueue:
		SetRepeatMode(RepeatModeTrack)
	case RepeatModeTrack:
		SetRepeatMode(RepeatModeNone)
	}
}

func Next() {
	logger.Debug("player controls requested to play next track")
	playNextTrack()
}

func Pause() {
	logger.Debug("player controls requested to pause")
	playbin.SetState(gst.StatePaused)
	PlaybackStateChanged.Notify(func(oldValue *PlaybackState) *PlaybackState {
		newState := *oldValue
		newState.Status = PlaybackStatusPaused
		return &newState
	})
}

func Play() {
	logger.Debug("player controls requested to start playing")
	playbin.SetState(gst.StatePlaying)
}

func PlayPause() {
	logger.Debug("player controls requested to start playing or pause")
	switch PlaybackStateChanged.CurrentValue().Status {
	case PlaybackStatusPlaying:
		Pause()
	case PlaybackStatusPaused:
		Play()
	case PlaybackStatusStopped:
		SeekToPosition(0)
	}
}

func Previous() {
	logger.Debug("player controls requested to play previous track")
	playPreviousTrack()
}

func SeekToPercent(percent float64) {
	if percent < 0 || percent > 100 {
		percent = math.Max(0, math.Min(100, percent))
		logger.Warn("percent must be between 0 and 100, clamping to nearest value", "percent", percent)
	}

	position := float64(PlaybackStateChanged.CurrentValue().Duration) * (percent / 100.0)
	SeekToPosition(time.Duration(int64(position)))
}

var (
	seekMutex              sync.Mutex
	seekDebounce           *time.Timer
	seekOriginalState      gst.State
	wasUpdateRunnerRunning bool
)

func SeekToPosition(position time.Duration) {
	logger.Debug("player controls requested to seek to position", "position", position)
	go func() {
		seekMutex.Lock()
		defer seekMutex.Unlock()

		if seekDebounce != nil {
			seekDebounce.Stop()
		} else {
			seekOriginalState = playbin.GetCurrentState()
			playbin.SetState(gst.StatePaused)

			if updateRunnerSourceHandle != 0 {
				wasUpdateRunnerRunning = true
				stopUpdateRunner()
			}
		}

		seekDebounce = time.AfterFunc(100*time.Millisecond, func() {
			seekMutex.Lock()
			defer seekMutex.Unlock()

			fmt.Println("Seeking. Playbin state ", seekOriginalState, " Update runner running", wasUpdateRunnerRunning)

			if PlaybackStateChanged.CurrentValue().Position == position {
				return
			}

			playbin.SeekTime(position, gst.SeekFlagFlush)
			playbin.SetState(seekOriginalState)

			if wasUpdateRunnerRunning {
				startUpdateRunner()
			}

			seekDebounce = nil

			PlaybackStateChanged.Notify(func(oldValue *PlaybackState) *PlaybackState {
				newState := *oldValue
				newState.Position = position
				newState.IsSeeking = true
				return &newState
			})
		})
	}()
}

func SeekToPositionRelative(delta time.Duration) {
	SeekToPosition(PlaybackStateChanged.CurrentValue().Position + delta)
}

func SetRepeatMode(m RepeatMode) {
	if m == RepeatModeChanged.CurrentValue() {
		return
	}

	RepeatModeChanged.Notify(func(oldValue RepeatMode) RepeatMode {
		return m
	})
}

func SetVolume(volume float64) {
	if volume < 0 {
		logger.Info("Volume is lower than 0, overriding back to 0.")
		volume = 0
	} else if volume > 1 {
		logger.Warn("Volume is higher than 1. This will cause overdrive to the speakers.")
	}
	playbin.SetProperty("volume", volume)
}

func Stop() {
	logger.Debug("player controls requested to stop")
	playbin.SetState(gst.StateNull)
}

func ToggleShuffle() {
	if ShuffleSeedChanged.CurrentValue() == 0 {
		ShuffleSeedChanged.Notify(func(oldValue int64) int64 {
			return time.Now().Unix()
		})
	} else {
		ShuffleSeedChanged.Notify(func(oldValue int64) int64 {
			return 0
		})
	}
}
