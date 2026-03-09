package v2

import (
	"strconv"

	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi"
	v2 "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/models/v2"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
)

var artistInfoLogger = logger.With("type", "ArtistInfo").WithGroup("artist_info")

type ArtistInfo struct {
	v2.ArtistItemData
}

func (a ArtistInfo) ID() string {
	return strconv.Itoa(a.Id)
}

func (a ArtistInfo) Cover(preferredSize int) string {
	logger := artistInfoLogger.With("method", "ProfilePicture").WithGroup("profile_picture").With("preferred_size", preferredSize)

	if preferredSize > 0 {
		logger.Debug("legacy api does not support preferred size")
	}

	if a.ArtistItemData.Picture == "" {
		return ""
	}

	return tidalapi.ImageURL(a.ArtistItemData.Picture)
}

func (a ArtistInfo) Route() string {
	return "artist/" + a.ID()
}

func (a ArtistInfo) SourceType() sonami.SourceType {
	return sonami.SourceTypeArtist
}

func (a ArtistInfo) Title() string {
	return a.ArtistItemData.Name
}

func (a ArtistInfo) URL() string {
	return "https://tidal.com/artist/" + a.ID()
}

func NewArtistInfo(artist v2.ArtistItemData) sonami.ArtistInfo {
	return &ArtistInfo{
		ArtistItemData: artist,
	}
}
