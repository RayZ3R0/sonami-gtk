package openapi

import (
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/models/openapi"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
)

var playlistLogger = logger.With("type", "Playlist").WithGroup("playlist")

type Playlist struct {
	PlaylistInfo
}

func (p Playlist) Count() int {
	return p.Data.Attributes.NumberOfItems
}

func (p Playlist) Creator() sonami.ArtistInfo {
	artists := p.Included.Artists(p.Data.Relationships.OwnerProfiles.Data...)

	if len(artists) == 0 {
		return nil
	}
	return NewArtistInfo(artists[0])
}

func NewPlaylist(playlist openapi.Playlist) sonami.Playlist {
	return &Playlist{PlaylistInfo{playlist}}
}
