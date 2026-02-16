package v2

import (
	"strconv"

	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	v2 "codeberg.org/dergs/tonearm/pkg/tidalapi/models/v2"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
)

var artistInfoLogger = logger.With("type", "ArtistInfo").WithGroup("artist_info")

type ArtistInfo struct {
	v2.ArtistItemData
}

func (a ArtistInfo) ID() string {
	return strconv.Itoa(a.Id)
}

func (a ArtistInfo) Name() string {
	return a.ArtistItemData.Name
}

func (a ArtistInfo) ProfilePicture(preferredSize int) string {
	logger := artistInfoLogger.With("method", "ProfilePicture").WithGroup("profile_picture").With("preferred_size", preferredSize)

	if preferredSize > 0 {
		logger.Debug("legacy api does not support preferred size")
	}

	if a.ArtistItemData.Picture == "" {
		return ""
	}

	return tidalapi.ImageURL(a.ArtistItemData.Picture)
}

func (a ArtistInfo) URL() string {
	return "https://tidal.com/artist/" + a.ID()
}

func NewArtistInfo(artist v2.ArtistItemData) tonearm.ArtistInfo {
	return &ArtistInfo{
		ArtistItemData: artist,
	}
}
