package openapi

import (
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
)

type MyTracksInfo struct{}

func (m MyTracksInfo) Cover(preferredSize int) string {
	return "https://tidal.com/assets/my-tracks-DTG3pLQW.png"
}

func (m MyTracksInfo) Route() string {
	return "my-collection/tracks"
}

func (m MyTracksInfo) Title() string {
	return gettext.Get("My Tracks")
}

func (m MyTracksInfo) SourceType() sonami.SourceType {
	return sonami.SourceTypePlaylist
}

func (m MyTracksInfo) URL() string {
	return "https://tidal.com/my-collection/tracks"
}
