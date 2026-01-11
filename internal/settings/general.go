package settings

import (
	"github.com/jwijenbergh/puregotk/v4/gio"
	"github.com/jwijenbergh/puregotk/v4/gobject"
)

type GeneralSettings struct {
	settings *gio.Settings
}

func (g *GeneralSettings) BindRunInBackground(target *gobject.Object, property string) {
	g.settings.Bind("allow-running-in-background", target, property, gio.GSettingsBindNoSensitivityValue)
}

func (g *GeneralSettings) BindDefaultPage(target *gobject.Object, property string) {
	g.settings.Bind("default-page", target, property, gio.GSettingsBindNoSensitivityValue)
}

func (g *GeneralSettings) ShouldRunInBackground() bool {
	return g.settings.GetBoolean("allow-running-in-background")
}

func (g *GeneralSettings) DefaultPage() string {
	return g.settings.GetString("default-page")
}

func (g *GeneralSettings) GetWindowHeight() int {
	return g.settings.GetInt("window-height")
}

func (g *GeneralSettings) GetWindowWidth() int {
	return g.settings.GetInt("window-width")
}

func (g *GeneralSettings) SetDefaultPage(path string) {
	g.settings.SetString("default-page", path)
}

func (g *GeneralSettings) SetWindowHeight(height int) {
	g.settings.SetInt("window-height", height)
}

func (g *GeneralSettings) SetWindowWidth(width int) {
	g.settings.SetInt("window-width", width)
}
