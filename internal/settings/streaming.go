package settings

import (
	"codeberg.org/puregotk/puregotk/v4/gio"
	"codeberg.org/puregotk/puregotk/v4/gobject"
)

type StreamingSettings struct {
	settings *gio.Settings
}

func (s *StreamingSettings) GetInstancesURL() string {
	return s.settings.GetString("instances-url")
}

func (s *StreamingSettings) SetInstancesURL(url string) {
	s.settings.SetString("instances-url", url)
}

func (s *StreamingSettings) BindInstancesURL(target *gobject.Object, property string) {
	s.settings.Bind("instances-url", target, property, gio.GSettingsBindNoSensitivityValue)
}
