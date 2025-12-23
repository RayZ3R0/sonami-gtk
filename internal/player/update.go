package player

import (
	"time"

	"github.com/go-gst/go-glib/glib"
	"github.com/go-gst/go-gst/gst"
)

var playbackStateUpdater glib.SourceHandle

func startUpdateRunner() {
	if playbackStateUpdater != 0 {
		return
	}

	playbackStateUpdater, _ = glib.TimeoutAdd(250, onUpdateTick, nil)
}

func stopUpdateRunner() {
	if playbackStateUpdater == 0 {
		return
	}

	glib.SourceRemove(playbackStateUpdater)
	playbackStateUpdater = 0
}

func onUpdateTick() bool {
	OnState.Notify(func(state *State) {
		if ok, duration := playbin.QueryDuration(gst.FormatTime); ok {
			state.Duration = int(time.Duration(duration).Seconds())
		}

		if ok, position := playbin.QueryPosition(gst.FormatTime); ok {
			state.Position = int(time.Duration(position).Seconds())
		}
	})
	return true
}

func onVolumeChange() {
	OnState.Notify(func(state *State) {
		if volume, err := playbin.GetProperty("volume"); err == nil {
			state.Volume = volume.(float64)
		}
	})
}
