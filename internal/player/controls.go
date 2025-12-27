package player

import (
	"time"

	"github.com/go-gst/go-gst/gst"
)

func Scrub(percent float64) {
	if percent < 0 || percent > 100 {
		panic("percent must be between 0 and 100")
	}

	position := float64(OnStateChanged.current.Duration) * percent / 100.0
	playbin.SeekTime(time.Duration(position)*time.Second, gst.SeekFlagFlush)
	OnStateChanged.Notify(func(state *State) {
		state.Position = int(position)
		if state.Status == StatusStopped {
			state.Status = StatusPlaying
		}
	})
}

func SeekTo(timestamp time.Duration) {
	playbin.SeekTime(timestamp, gst.SeekFlagFlush)
	OnStateChanged.Notify(func(state *State) {
		state.Position = int(timestamp.Seconds())
		if state.Status == StatusStopped {
			state.Status = StatusPlaying
		}
	})
}

func PlayPause() {
	if OnStateChanged.current.Status == StatusPlaying {
		playbin.SetState(gst.StatePaused)
		OnStateChanged.Notify(func(state *State) {
			state.Status = StatusPaused
		})
	} else if OnStateChanged.current.Status == StatusPaused {
		playbin.SetState(gst.StatePlaying)
		OnStateChanged.Notify(func(state *State) {
			state.Status = StatusPlaying
		})
	} else if OnStateChanged.current.Status == StatusStopped {
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
