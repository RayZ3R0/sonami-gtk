package openapi

import (
	"strings"
	"time"

	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/models/openapi"
)

var playlistInfoLogger = logger.With("type", "PlaylistInfo").WithGroup("playlist_info")

type PlaylistInfo struct {
	openapi.Playlist
}

func (p PlaylistInfo) Cover(preferredSize int) string {
	logger := playlistInfoLogger.With("method", "Cover").WithGroup("cover").With("preferred_size", preferredSize)

	if preferredSize < 0 {
		logger.Debug("defaulting to smallest picture size as a negative value was passed")
		preferredSize = 0
	}

	artworks := p.Included.PlainArtworks(p.Data.Relationships.CoverArt.Data...)

	return artworks.AtLeast(preferredSize)
}

func (p PlaylistInfo) CreatedAt() time.Time {
	return p.Data.Attributes.CreatedAt.Time
}

func (p PlaylistInfo) Duration() time.Duration {
	if p.Data.Attributes.Duration == nil {
		return time.Duration(0)
	}
	return p.Data.Attributes.Duration.Duration
}

func (p PlaylistInfo) ID() string {
	return p.Playlist.Data.ID
}

func (p PlaylistInfo) IsMix() bool {
	return !strings.Contains(p.ID(), "-")
}

func (p PlaylistInfo) Route() string {
	return "playlist/" + p.ID()
}

func (p PlaylistInfo) SourceType() sonami.SourceType {
	return sonami.SourceTypePlaylist
}

func (p PlaylistInfo) Title() string {
	return p.Playlist.Data.Attributes.Name
}

func (p PlaylistInfo) URL() string {
	if p.IsMix() {
		return "https://tidal.com/mix/" + p.ID()
	}
	return "https://tidal.com/playlist/" + p.ID()
}

func NewPlaylistInfo(playlist openapi.Playlist) sonami.PlaylistInfo {
	return &PlaylistInfo{playlist}
}
