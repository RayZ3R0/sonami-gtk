package player

import (
	"time"

	"github.com/go-gst/go-gst/gst"
)

func Scrub(percent float64) {
	if percent < 0 || percent > 100 {
		panic("percent must be between 0 and 100")
	}

	position := time.Duration(int64(float64(OnStateChanged.current.Duration.Nanoseconds()) * (percent / 100.0)))
	playbin.SeekTime(position, gst.SeekFlagFlush)
	OnStateChanged.Notify(func(state *State) {
		state.Position = position
		if state.Status == StatusStopped {
			state.Status = StatusPlaying
		}
	})
}

// Removes duration to the current position
func SeekBackward(duration time.Duration) {
	SeekForward(-duration)
}

// Adds duration to the current position
func SeekForward(duration time.Duration) {
	ok, position := playbin.QueryPosition(gst.FormatTime)
	if !ok {
		return
	}
	SeekTo(time.Duration(position) + duration)
}

func SeekTo(timestamp time.Duration) {
	playbin.SeekTime(timestamp, gst.SeekFlagFlush)
	OnStateChanged.Notify(func(state *State) {
		state.Position = timestamp
		if state.Status == StatusStopped {
			state.Status = StatusPlaying
		}
	})
}

func Pause() {
	playbin.SetState(gst.StatePaused)
	OnStateChanged.Notify(func(state *State) {
		state.Status = StatusPaused
	})
}

func Play() {
	playbin.SetState(gst.StatePlaying)
	OnStateChanged.Notify(func(state *State) {
		state.Status = StatusPlaying
	})
}

func PlayPause() {
	switch OnStateChanged.current.Status {
	case StatusPlaying:
		Pause()
	case StatusPaused:
		Play()
	case StatusStopped:
		Scrub(0)
	}
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

func Next() {
	nextTrack()
}

func Previous() {
	ok, position := playbin.QueryPosition(gst.FormatTime)
	if ok && time.Duration(position) > 5*time.Second {
		SeekTo(0)
		return
	}

	if len(history.Entries) < 1 {
		SeekTo(0)
		return
	}
	entry := history.Pop()
	if entry != nil {
		// Re-Queue current song to front of user-queue
		UserQueue.AddTrackID(OnTrackChanged.current.ID, true)
		// Switch to previous track without clearing base queue
		playTrackId(entry.TrackID)
	}
}
