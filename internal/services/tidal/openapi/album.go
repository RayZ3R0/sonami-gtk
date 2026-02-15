package openapi

import (
	"context"

	"codeberg.org/dergs/tonearm/internal/services/tidal/internal"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
)

var albumLogger = logger.With("type", "album")

type Album struct {
	API   *tidalapi.TidalAPI
	Album openapi.Album
}

func (a *Album) Artists(hints ...tonearm.FetchHint) (tonearm.Paginator[tonearm.Artist], error) {
	if len(a.Album.Data.Relationships.Artists.Data) > 0 {
		return internal.NewPaginatorWithFirstPage(internal.PaginatedOpenAPIFunc(a.API.OpenAPI.V2.Albums.Artists), a.ID(), func(r *openapi.Response[[]openapi.Relationship]) []tonearm.Artist {
			return ArtistsFromIncluded(a.API, r.Included, r.Data...)
		}, a.Album.Data.Relationships.Artists.Links.Next, ArtistsFromIncluded(a.API, a.Album.Included, a.Album.Data.Relationships.Artists.Data...), ArtistHintsToInclude("", hints...)...), nil
	} else {
		return internal.NewPaginator(internal.PaginatedOpenAPIFunc(a.API.OpenAPI.V2.Albums.Artists), a.ID(), func(r *openapi.Response[[]openapi.Relationship]) []tonearm.Artist {
			return ArtistsFromIncluded(a.API, r.Included, r.Data...)
		}, ArtistHintsToInclude("", hints...)...), nil
	}
}

func (a *Album) Cover(preferredSize int) (string, error) {
	albumCoverArtLogger := albumLogger.With("method", "Cover")

	var artworks openapi.Artworks
	if len(a.Album.Data.Relationships.CoverArt.Data) > 0 {
		albumCoverArtLogger.Debug("using pre-fetched cover artworks")
		artworks = a.Album.Included.PlainArtworks(a.Album.Data.Relationships.CoverArt.Data...)
	} else {
		albumCoverArtLogger.Debug("requesting cover artworks (are we missing a hint when fetching the album?)")
		album, err := a.API.OpenAPI.V2.Albums.Album(context.Background(), a.ID(), "coverArt")
		if err != nil {
			return "", err
		}
		artworks = album.Included.PlainArtworks(album.Data.Relationships.CoverArt.Data...)
	}

	return artworks.AtLeast(preferredSize), nil
}

func (a *Album) ID() string {
	return a.Album.Data.ID
}

func (a *Album) URL() string {
	return "https://tidal.com/album/" + a.ID()
}

func (a *Album) Route() string {
	return "album/" + a.ID()
}

func (a *Album) Title() string {
	return a.Album.Data.Attributes.Title
}

func (a *Album) Tracks(hints ...tonearm.FetchHint) (tonearm.Paginator[tonearm.Track], error) {
	if len(a.Album.Data.Relationships.Items.Data) > 0 {
		return internal.NewPaginatorWithFirstPage(internal.PaginatedOpenAPIFunc(a.API.OpenAPI.V2.Albums.Items), a.ID(), func(r *openapi.Response[[]openapi.Relationship]) []tonearm.Track {
			return TracksFromIncluded(a.API, r.Included, r.Data...)
		}, a.Album.Data.Relationships.Items.Links.Next, TracksFromIncluded(a.API, a.Album.Included, a.Album.Data.Relationships.Items.Data...)), nil
	} else {
		return internal.NewPaginator(internal.PaginatedOpenAPIFunc(a.API.OpenAPI.V2.Albums.Items), a.ID(), func(r *openapi.Response[[]openapi.Relationship]) []tonearm.Track {
			return TracksFromIncluded(a.API, r.Included, r.Data...)
		}), nil
	}
}

func AlbumsFromIncluded(api *tidalapi.TidalAPI, included openapi.IncludedObjects, relationships ...openapi.Relationship) []tonearm.Album {
	var albums []tonearm.Album
	for _, album := range included.Albums(relationships...) {
		albums = append(albums, NewAlbum(api, album))
	}
	return albums
}

func NewAlbum(api *tidalapi.TidalAPI, album openapi.Album) tonearm.Album {
	return &Album{
		API:   api,
		Album: album,
	}
}

func AlbumHintsToInclude(prefix string, hints ...tonearm.FetchHint) []string {
	included := []string{}
	for _, hint := range hints {
		switch hint {
		case tonearm.AlbumHintCover:
			included = append(included, prefix+"coverArt")
		case tonearm.AlbumHintArtists:
			included = append(included, prefix+"artists")
		case tonearm.AlbumHintTracks:
			included = append(included, prefix+"items")
		}
	}
	return included
}
