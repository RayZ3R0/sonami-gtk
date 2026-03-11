package settings

import (
	"codeberg.org/puregotk/puregotk/v4/gio"
	"codeberg.org/puregotk/puregotk/v4/gobject"
)

type DiscordSettings struct {
	settings *gio.Settings
}

func (d *DiscordSettings) RichPresenceEnabled() bool {
	return d.settings.GetBoolean("enable-rich-presence")
}

func (d *DiscordSettings) BindEnableRichPresence(target *gobject.Object, property string) {
	d.settings.Bind("enable-rich-presence", target, property, gio.GSettingsBindNoSensitivityValue)
}
