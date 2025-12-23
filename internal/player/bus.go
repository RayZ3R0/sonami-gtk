package player

import (
	"log/slog"

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
	OnState.Notify(func(state *State) {
		state.Status = StatusStopped
		state.Position = 0
		state.Duration = 0
		state.Track = nil
	})
}
