package v1

import (
	"fmt"

	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
	v1 "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/models/v1"
)

type Track struct {
	TrackInfo
	album sonami.AlbumInfo
}

func (t Track) Album() sonami.AlbumInfo {
	return t.album
}

func (t Track) Artists() sonami.ArtistInfos {
	artists := make(sonami.ArtistInfos, len(t.AlbumItem.Artists))
	for i, artist := range t.AlbumItem.Artists {
		artists[i] = NewArtistInfo(artist)
	}
	return artists
}

func (t Track) Cover(perferredSize int) string {
	return t.Album().Cover(perferredSize)
}

func (t Track) Route() string {
	return fmt.Sprintf("album/%s", t.Album().ID())
}

func (t Track) SourceType() sonami.SourceType {
	return sonami.SourceTypeTrack
}

func NewTrack(item v1.AlbumItem, album sonami.AlbumInfo) sonami.Track {
	return &Track{TrackInfo{AlbumItem: item}, album}
}
