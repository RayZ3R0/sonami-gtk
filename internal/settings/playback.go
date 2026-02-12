package settings

import (
	"github.com/jwijenbergh/puregotk/v4/gio"
	"github.com/jwijenbergh/puregotk/v4/gobject"
)

type ReplayGainMode string

const (
	ReplayGainModeAuto  ReplayGainMode = "Auto"
	ReplayGainModeAlbum ReplayGainMode = "Album"
	ReplayGainModeTrack ReplayGainMode = "Track"
)

var ReplayGainModes = []ReplayGainMode{
	ReplayGainModeAuto,
	ReplayGainModeAlbum,
	ReplayGainModeTrack,
}

var ReplayGainModeStrings = make([]string, len(ReplayGainModes))

var ReplayGainModeIndex = make(map[ReplayGainMode]uint)

func init() {
	for index, mode := range ReplayGainModes {
		ReplayGainModeIndex[mode] = uint(index)
		ReplayGainModeStrings[index] = string(mode)
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
	return ReplayGainMode(p.settings.GetString("replay-gain-mode"))
}

func (p *PlaybackSettings) ReplayGainModeIndex() uint {
	return ReplayGainModeIndex[p.ReplayGainMode()]
}

func (p *PlaybackSettings) SetReplayGainMode(mode ReplayGainMode) {
	p.settings.SetString("replay-gain-mode", string(mode))
}

func (p *PlaybackSettings) SetReplayGainModeIndex(index uint) {
	p.SetReplayGainMode(ReplayGainModes[index])
}
