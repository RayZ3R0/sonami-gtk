package v2

import (
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
	v2 "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/models/v2"
)

var playlistLogger = logger.With("type", "Playlist").WithGroup("playlist")

type Playlist struct {
	PlaylistInfo
}

func (p Playlist) Count() int {
	return p.PlaylistItemData.NumberOfTracks
}

func (p Playlist) Creator() sonami.ArtistInfo {
	if p.PlaylistItemData.Creator.ID == 0 {
		return nil
	}

	return NewArtistInfo(v2.ArtistItemData{
		Id:      p.PlaylistItemData.Creator.ID,
		Name:    p.PlaylistItemData.Creator.Name,
		Picture: p.PlaylistItemData.Creator.Picture,
	})
}

func NewPlaylist(playlist v2.PlaylistItemData) sonami.Playlist {
	return &Playlist{PlaylistInfo{playlist}}
}
