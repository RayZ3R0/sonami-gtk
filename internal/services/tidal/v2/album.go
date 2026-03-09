package v2

import (
	v2 "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/models/v2"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
)

var albumLogger = logger.With("type", "Album").WithGroup("album")

type Album struct {
	AlbumInfo
}

func (a *Album) Artists() sonami.ArtistInfos {
	artists := make(sonami.ArtistInfos, 0)
	for _, artist := range a.AlbumItemData.Artists {
		artists = append(artists, NewArtistInfo(artist))
	}
	return artists
}

func (a AlbumInfo) Count() int {
	logger := albumLogger.With("method", "Count").WithGroup("count")
	logger.Warn("v2 album does not support Count")
	return -1
}

func NewAlbum(album v2.AlbumItemData) sonami.Album {
	return &Album{AlbumInfo{album}}
}
