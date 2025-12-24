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
	logger.Debug("started playbin update runner", "source_handle", playbackStateUpdater)
}

func stopUpdateRunner() {
	if playbackStateUpdater == 0 {
		return
	}

	glib.SourceRemove(playbackStateUpdater)
	logger.Debug("stopped playbin update runner", "source_handle", playbackStateUpdater)
	playbackStateUpdater = 0
}

func onUpdateTick() bool {
	OnStateChanged.Notify(func(state *State) {
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
	OnVolumeChanged.Notify(func(previous float64) float64 {
		if volume, err := playbin.GetProperty("volume"); err == nil {
			return volume.(float64)
		}
		return previous
	})
}
