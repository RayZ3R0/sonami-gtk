package settings

import (
	"github.com/jwijenbergh/puregotk/v4/gio"
	"github.com/jwijenbergh/puregotk/v4/gobject"
)

type ScrobblingSettings struct {
	settings *gio.Settings
}

func (s *ScrobblingSettings) BindEnableListenBrainz(target *gobject.Object, property string) {
	s.settings.Bind("enable-listenbrainz", target, property, gio.GSettingsBindNoSensitivityValue)
}

func (s *ScrobblingSettings) BindListenBrainzToken(target *gobject.Object, property string) {
	s.settings.Bind("listenbrainz-token", target, property, gio.GSettingsBindNoSensitivityValue)
}

func (s *ScrobblingSettings) ShouldEnableListenBrainz() bool {
	return s.settings.GetBoolean("enable-listenbrainz")
}
func (s *ScrobblingSettings) ListenBrainzToken() string {
	return s.settings.GetString("listenbrainz-token")
}
