package openapi

import (
	"context"
	"fmt"

	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
)

var artistLogger = logger.With("type", "artist")

type Artist struct {
	API    *tidalapi.TidalAPI
	Artist openapi.Artist
}

func (a *Artist) ID() string {
	return a.Artist.Data.ID
}

func (a *Artist) Name() string {
	return a.Artist.Data.Attributes.Name
}

func (a *Artist) ProfilePicture(preferredSize int) (string, error) {
	logger := artistLogger.With("method", "ProfilePicture")

	var artworks openapi.Artworks
	if len(a.Artist.Data.Relationships.ProfileArt.Data) > 0 {
		logger.Debug("using pre-fetched profile artworks")
		artworks = a.Artist.Included.PlainArtworks(a.Artist.Data.Relationships.ProfileArt.Data...)
	} else {
		logger.Debug("requesting profile artworks (are we missing a hint when fetching the artist?)")
		album, err := a.API.OpenAPI.V2.Artists.Artist(context.Background(), a.ID(), "profileArt")
		if err != nil {
			return "", err
		}
		artworks = album.Included.PlainArtworks(album.Data.Relationships.ProfileArt.Data...)
	}

	return artworks.AtLeast(preferredSize), nil
}

func (a *Artist) URL() string {
	return fmt.Sprintf("https://tidal.com/artist/%s", a.ID())
}

func ArtistsFromIncluded(api *tidalapi.TidalAPI, included openapi.IncludedObjects, relationships ...openapi.Relationship) []tonearm.Artist {
	var artists []tonearm.Artist
	for _, artist := range included.Artists(relationships...) {
		artists = append(artists, NewArtist(api, artist))
	}
	return artists
}

func NewArtist(api *tidalapi.TidalAPI, artist openapi.Artist) tonearm.Artist {
	return &Artist{
		API:    api,
		Artist: artist,
	}
}

func ArtistHintsToInclude(prefix string, hints ...tonearm.FetchHint) []string {
	included := []string{}
	for _, hint := range hints {
		switch hint {
		case tonearm.ArtistHintProfilePicture:
			included = append(included, prefix+"profileArt")
		}
	}
	return included
}
