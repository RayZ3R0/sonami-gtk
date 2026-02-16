package v2

import (
	v2 "codeberg.org/dergs/tonearm/pkg/tidalapi/models/v2"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
)

var trackLogger = logger.With("type", "Track").WithGroup("track")

type Track struct {
	TrackInfo
}

func (t Track) Album() tonearm.AlbumInfo {
	logger := trackLogger.With("method", "Album").WithGroup("album")
	logger.Debug("v2 track does not properly resolve album artists or duration")
	return NewAlbumInfo(v2.AlbumItemData{
		Artists:     t.TrackItemData.Artists,
		Cover:       t.TrackItemData.Album.Cover,
		Id:          t.TrackItemData.Album.ID,
		Duration:    0,
		ReleaseDate: t.TrackItemData.Album.ReleaseDate,
		Title:       t.TrackItemData.Album.Title,
		Type:        "ALBUM",
	})
}

func (t Track) Artists() tonearm.ArtistInfos {
	artists := make(tonearm.ArtistInfos, len(t.TrackItemData.Artists))
	for i, artist := range t.TrackItemData.Artists {
		artists[i] = NewArtistInfo(artist)
	}
	return artists
}

func (t Track) Cover(perferredSize int) string {
	return t.Album().Cover(perferredSize)
}

func NewTrack(item v2.TrackItemData) tonearm.Track {
	return &Track{TrackInfo{item}}
}
