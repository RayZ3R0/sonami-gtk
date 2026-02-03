package callback

import (
	"codeberg.org/dergs/tonearm/pkg/utils/cutil"
	"github.com/jwijenbergh/puregotk/v4/gio"
)

var (
	GioSettingsChangedCallback = func(settings gio.Settings, setting string) {
		CallbackHandler[any](settings.Object, "changed", settings, cutil.ParseNullTerminatedString(setting))
	}
)
