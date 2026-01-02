package callback

import "github.com/jwijenbergh/puregotk/v4/gtk"

var (
	SearchEntryActivateCallback = func(widget gtk.SearchEntry) {
		CallbackHandler[any](widget.Widget, "activate", widget)
	}
	SearchChangedCallback = func(widget gtk.SearchEntry) {
		CallbackHandler[any](widget.Widget, "search-changed", widget)
	}
)
