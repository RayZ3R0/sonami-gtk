package v1

import (
	"fmt"

	v1 "codeberg.org/dergs/tonearm/pkg/tidalapi/models/v1"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
)

type Track struct {
	TrackInfo
	album tonearm.AlbumInfo
}

func (t Track) Album() tonearm.AlbumInfo {
	return t.album
}

func (t Track) Artists() tonearm.ArtistInfos {
	artists := make(tonearm.ArtistInfos, len(t.AlbumItem.Artists))
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

func (t Track) SourceType() tonearm.SourceType {
	return tonearm.SourceTypeTrack
}

func NewTrack(item v1.AlbumItem, album tonearm.AlbumInfo) tonearm.Track {
	return &Track{TrackInfo{AlbumItem: item}, album}
}
