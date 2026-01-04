package callback

import "github.com/jwijenbergh/puregotk/v4/gtk"

var (
	ScrolledWindowEdgeReachedCallback = func(widget gtk.ScrolledWindow, positionType gtk.PositionType) {
		CallbackHandler[any](widget.Widget, "edge-reached", widget, positionType)
	}
)
