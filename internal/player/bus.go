package player

import (
	"log/slog"
	"strconv"

	"github.com/go-gst/go-gst/gst"
)

func onBusMessage(msg *gst.Message) bool {
	switch msg.Type() {
	case gst.MessageError:
		onBusError(msg.ParseError())
	case gst.MessageStateChanged:
		onBusStateChanged(msg.ParseStateChanged())
	case gst.MessageEOS:
		onBusStreamEnd()
	case gst.MessageStreamStart:
		onBusStreamStart()
	}
	return true
}

func onBusError(err *gst.GError) {
	slog.Error("Error while playing track", "error", err.Error())
}

func onBusStateChanged(_, newState gst.State) {
	if newState == gst.StatePlaying {
		startUpdateRunner()
	}
}

func onBusStreamEnd() {
	stopUpdateRunner()
	OnStateChanged.Notify(func(state *State) {
		state.Status = StatusStopped
	})
}

func onBusStreamStart() {
	// A hack to trigger the correct track updates with gapless playback
	if OnTrackChanged.current.ID != strconv.Itoa(currentlyConfiguredTrackID) {
		nextTrack()
	}
}
