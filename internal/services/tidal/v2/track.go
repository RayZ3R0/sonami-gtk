package v2

import (
	"fmt"

	v2 "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/models/v2"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
)

var trackLogger = logger.With("type", "Track").WithGroup("track")

type Track struct {
	TrackInfo
}

func (t Track) Album() sonami.AlbumInfo {
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

func (t Track) Artists() sonami.ArtistInfos {
	artists := make(sonami.ArtistInfos, len(t.TrackItemData.Artists))
	for i, artist := range t.TrackItemData.Artists {
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

func NewTrack(item v2.TrackItemData) sonami.Track {
	return &Track{TrackInfo{item}}
}
