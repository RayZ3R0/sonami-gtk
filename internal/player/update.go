package player

import (
	"time"

	"github.com/go-gst/go-gst/gst"
	"github.com/jwijenbergh/puregotk/v4/glib"
)

var playbackStateUpdater uint

const UpdateInterval = 250 * time.Millisecond

func startUpdateRunner() {
	if playbackStateUpdater != 0 {
		return
	}

	playbackStateUpdater = glib.TimeoutAdd(uint(UpdateInterval.Milliseconds()), &onUpdateTick, 0)
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

var onUpdateTick = glib.SourceFunc(func(uintptr) bool {
	OnStateChanged.Notify(func(state *State) {
		if ok, duration := playbin.QueryDuration(gst.FormatTime); ok {
			state.Duration = time.Duration(duration)
		}

		if ok, position := playbin.QueryPosition(gst.FormatTime); ok {
			state.Position = time.Duration(position)
		}
	})
	return glib.SOURCE_CONTINUE
})

func onVolumeChange() {
	OnVolumeChanged.Notify(func(previous float64) float64 {
		if volume, err := playbin.GetProperty("volume"); err == nil {
			return volume.(float64)
		}
		return previous
	})
}
