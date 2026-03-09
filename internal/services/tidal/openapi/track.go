package openapi

import (
	"fmt"

	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/models/openapi"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
)

var trackLogger = logger.With("type", "Track").WithGroup("track")

type Track struct {
	TrackInfo
}

func (t Track) Album() sonami.AlbumInfo {
	albums := t.Included.Albums(t.Data.Relationships.Albums.Data...)

	for _, album := range albums {
		return NewAlbumInfo(album)
	}
	return nil
}

func (t Track) Artists() sonami.ArtistInfos {
	artworks := t.Included.Artists(t.Data.Relationships.Artists.Data...)

	artists := make(sonami.ArtistInfos, 0)
	for _, artist := range artworks {
		artists = append(artists, NewArtistInfo(artist))
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

func NewTrack(item openapi.Track) sonami.Track {
	return &Track{TrackInfo{item}}
}
