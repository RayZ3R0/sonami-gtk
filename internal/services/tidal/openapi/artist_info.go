package openapi

import (
	"codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
)

var artistInfoLogger = logger.With("type", "ArtistInfo").WithGroup("artist_info")

type ArtistInfo struct {
	openapi.Artist
}

func (a ArtistInfo) ID() string {
	return a.Artist.Data.ID
}

func (a ArtistInfo) Name() string {
	return a.Artist.Data.Attributes.Name
}

func (a ArtistInfo) ProfilePicture(preferredSize int) string {
	logger := artistInfoLogger.With("method", "ProfilePicture").WithGroup("profile_picture").With("preferred_size", preferredSize)

	if preferredSize < 0 {
		logger.Debug("defaulting to smallest picture size as a negative value was passed")
		preferredSize = 0
	}

	artworks := a.Included.PlainArtworks(a.Data.Relationships.ProfileArt.Data...)
	logger.Debug("resolved profile artworks", "count", len(artworks))

	return artworks.AtLeast(preferredSize)
}

func (a ArtistInfo) URL() string {
	return "https://tidal.com/artist/" + a.ID()
}

func NewArtistInfo(artist openapi.Artist) tonearm.ArtistInfo {
	return ArtistInfo{artist}
}
