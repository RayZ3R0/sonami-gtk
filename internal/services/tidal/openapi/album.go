package openapi

import (
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/models/openapi"
)

var albumLogger = logger.With("type", "Album").WithGroup("album")

type Album struct {
	AlbumInfo
}

func (a *Album) Artists() sonami.ArtistInfos {
	artworks := a.Included.Artists(a.Data.Relationships.Artists.Data...)

	artists := make(sonami.ArtistInfos, 0)
	for _, artist := range artworks {
		artists = append(artists, NewArtistInfo(artist))
	}
	return artists
}

func (a AlbumInfo) Count() int {
	return a.Data.Attributes.NumberOfItems
}

func NewAlbum(album openapi.Album) sonami.Album {
	return &Album{AlbumInfo{album}}
}
