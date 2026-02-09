package settings

import (
	"github.com/jwijenbergh/puregotk/v4/gio"
	"github.com/jwijenbergh/puregotk/v4/gobject"
)

type PlaybackSettings struct {
	settings *gio.Settings
}

func (p *PlaybackSettings) BindAllowAutoplay(target *gobject.Object, property string) {
	p.settings.Bind("allow-autoplay", target, property, gio.GSettingsBindNoSensitivityValue)
}

func (p *PlaybackSettings) BindNormalizeVolume(target *gobject.Object, property string) {
	p.settings.Bind("normalize-volume", target, property, gio.GSettingsBindNoSensitivityValue)
}

func (p *PlaybackSettings) AllowAutoplay() bool {
	return p.settings.GetBoolean("allow-autoplay")
}

func (p *PlaybackSettings) NormalizeVolume() bool {
	return p.settings.GetBoolean("normalize-volume")
}
