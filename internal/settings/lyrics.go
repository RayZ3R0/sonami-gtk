package settings

import (
	"codeberg.org/puregotk/puregotk/v4/gio"
)

const (
	LyricsProviderTidal   = "tidal"
	LyricsProviderNetEase = "netease"
)

var LyricsProviderStrings = []string{"TIDAL", "NetEase"}

type LyricsSettings struct {
	settings *gio.Settings
}

func (l *LyricsSettings) Provider() string {
	return l.settings.GetString("provider")
}

func (l *LyricsSettings) ProviderIndex() uint32 {
	if l.Provider() == LyricsProviderNetEase {
		return 1
	}
	return 0
}

func (l *LyricsSettings) SetProviderByIndex(i uint32) {
	switch i {
	case 1:
		l.settings.SetString("provider", LyricsProviderNetEase)
	default:
		l.settings.SetString("provider", LyricsProviderTidal)
	}
}
