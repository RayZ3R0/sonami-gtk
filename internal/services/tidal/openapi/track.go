package openapi

import (
	"codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
)

var trackLogger = logger.With("type", "Track").WithGroup("track")

type Track struct {
	TrackInfo
}

func (t Track) Album() tonearm.AlbumInfo {
	logger := trackLogger.With("method", "Album").WithGroup("album")

	albums := t.Included.Albums(t.Data.Relationships.Albums.Data...)
	logger.Debug("resolved track albums", "count", len(albums))

	for _, album := range albums {
		return NewAlbumInfo(album)
	}
	return nil
}

func (t Track) Artists() tonearm.ArtistInfos {
	logger := trackLogger.With("method", "Artists").WithGroup("artists")

	artworks := t.Included.Artists(t.Data.Relationships.Artists.Data...)
	logger.Debug("resolved track artists", "count", len(artworks))

	artists := make(tonearm.ArtistInfos, 0)
	for _, artist := range artworks {
		artists = append(artists, NewArtistInfo(artist))
	}
	return artists
}

func (t Track) Cover(perferredSize int) string {
	return t.Album().Cover(perferredSize)
}

func NewTrack(item openapi.Track) tonearm.Track {
	return &Track{TrackInfo{item}}
}
