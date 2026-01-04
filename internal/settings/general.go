package settings

import "github.com/jwijenbergh/puregotk/v4/gio"

type GeneralSettings struct {
	settings *gio.Settings
}

func (g *GeneralSettings) GetWindowHeight() int {
	return g.settings.GetInt("window-height")
}

func (g *GeneralSettings) GetWindowWidth() int {
	return g.settings.GetInt("window-width")
}

func (g *GeneralSettings) SetWindowHeight(height int) {
	g.settings.SetInt("window-height", height)
}

func (g *GeneralSettings) SetWindowWidth(width int) {
	g.settings.SetInt("window-width", width)
}
