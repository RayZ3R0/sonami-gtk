package openapi

import (
	"context"
	"errors"
	"time"

	"codeberg.org/dergs/tonearm/internal/services/tidal/internal"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
)

var trackLogger = logger.With("type", "artist")

type Track struct {
	API   *tidalapi.TidalAPI
	Track openapi.Track
}

func (t *Track) Album(hints ...tonearm.FetchHint) (tonearm.Album, error) {
	logger := trackLogger.With("method", "Album")

	var albums []openapi.Album
	if len(t.Track.Data.Relationships.Albums.Data) > 0 {
		logger.Debug("using pre-fetched albums")
		albums = t.Track.Included.Albums(t.Track.Data.Relationships.Albums.Data...)
	} else {
		logger.Debug("requesting albums (are we missing a hint when fetching the track?)")
		response, err := t.API.OpenAPI.V2.Tracks.Albums(context.Background(), t.ID(), "", append(AlbumHintsToInclude("", hints...), "albums")...)
		if err != nil {
			return nil, err
		}
		albums = response.Included.Albums(response.Data...)
	}

	if len(albums) == 0 {
		return nil, errors.New("no album found")
	}
	return NewAlbum(t.API, albums[0]), nil
}

func (t *Track) ArtistNames() ([]string, error) {
	artistPaginator, err := t.Artists()
	if err != nil {
		return nil, err
	}

	artists, err := artistPaginator.GetAll()
	if err != nil {
		return nil, err
	}

	names := make([]string, len(artists))
	for i, artist := range artists {
		names[i] = artist.Name()
	}
	return names, nil
}

func (t *Track) Artists(hints ...tonearm.FetchHint) (tonearm.Paginator[tonearm.Artist], error) {
	if len(t.Track.Data.Relationships.Artists.Data) > 0 {
		return internal.NewPaginatorWithFirstPage(internal.PaginatedOpenAPIFunc(t.API.OpenAPI.V2.Tracks.Artists), t.ID(), func(r *openapi.Response[[]openapi.Relationship]) []tonearm.Artist {
			return ArtistsFromIncluded(t.API, r.Included, r.Data...)
		}, t.Track.Data.Relationships.Artists.Links.Next, ArtistsFromIncluded(t.API, t.Track.Included, t.Track.Data.Relationships.Artists.Data...), ArtistHintsToInclude("", hints...)...), nil
	} else {
		return internal.NewPaginator(internal.PaginatedOpenAPIFunc(t.API.OpenAPI.V2.Tracks.Artists), t.ID(), func(r *openapi.Response[[]openapi.Relationship]) []tonearm.Artist {
			return ArtistsFromIncluded(t.API, r.Included, r.Data...)
		}, append(ArtistHintsToInclude("", hints...), "artists")...), nil
	}
}

func (t *Track) Cover(preferredSize int) (string, error) {
	album, err := t.Album()
	if err != nil {
		return "", err
	}
	return album.Cover(preferredSize)
}

func (t *Track) ID() string {
	return t.Track.Data.ID
}

func (t *Track) Duration() time.Duration {
	return t.Track.Data.Attributes.Duration.Duration
}

func (t *Track) Title() string {
	return t.Track.Data.Attributes.Title
}

func (t *Track) URL() string {
	return "https://tidal.com/track/" + t.Track.Data.ID
}

func TracksFromIncluded(api *tidalapi.TidalAPI, included openapi.IncludedObjects, relationships ...openapi.Relationship) []tonearm.Track {
	var tracks []tonearm.Track
	for _, track := range included.Tracks(relationships...) {
		tracks = append(tracks, NewTrack(api, track))
	}
	return tracks
}

func NewTrack(api *tidalapi.TidalAPI, track openapi.Track) tonearm.Track {
	return &Track{
		API:   api,
		Track: track,
	}
}

func TrackHintsToInclude(prefix string, hints ...tonearm.FetchHint) []string {
	included := []string{}
	for _, hint := range hints {
		switch hint {
		case tonearm.TrackHintAlbum:
			included = append(included, prefix+"albums")
		case tonearm.TrackHintArtists:
			included = append(included, prefix+"artists")
		}
	}
	return included
}
