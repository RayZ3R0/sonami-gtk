package settings

import (
	"codeberg.org/dergs/tonearm/internal/gettext"
	"github.com/jwijenbergh/puregotk/v4/gio"
	"github.com/jwijenbergh/puregotk/v4/gobject"
)

type ReplayGainMode int

const (
	// NEVER change the values of existing modes
	ReplayGainModeAuto  ReplayGainMode = 0
	ReplayGainModeAlbum ReplayGainMode = 1
	ReplayGainModeTrack ReplayGainMode = 2
)

func ReplayGainModeStrings() []string {
	return []string{
		gettext.Get("Auto"),
		gettext.Get("Album"),
		gettext.Get("Track"),
	}
}

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

func (p *PlaybackSettings) ReplayGainMode() ReplayGainMode {
	return ReplayGainMode(p.settings.GetInt("replay-gain-mode"))
}

func (p *PlaybackSettings) SetReplayGainMode(mode ReplayGainMode) {
	p.settings.SetInt("replay-gain-mode", int(mode))
}
