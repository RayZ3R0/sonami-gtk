package tidal

import (
	"context"
	"log/slog"

	"codeberg.org/dergs/tonearm/internal/secrets"
	"codeberg.org/dergs/tonearm/internal/services/tidal/openapi"
	v1 "codeberg.org/dergs/tonearm/internal/services/tidal/v1"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	modelopenapi "codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"
	modelv1 "codeberg.org/dergs/tonearm/pkg/tidalapi/models/v1"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
	"github.com/infinytum/injector"
)

var logger = slog.With("service", "TIDAL")

var (
	getAlbumLogger  = logger.With("method", "GetAlbum")
	getArtistLogger = logger.With("method", "GetArtist")
	getTrackLogger  = logger.With("method", "GetTrack")
)

func init() {
	injector.DeferredSingleton(func() *tidalapi.TidalAPI {
		countryCode, err := tidalapi.FetchCountryCode()
		if err != nil {
			slog.Error("Failed to fetch country code, defaulting to WW", err)
			countryCode = "WW"
		}
		slog.Info("Discovered country code", "countryCode", countryCode)
		return tidalapi.NewClient(countryCode, secrets.NewTokenAuthStrategy())
	})
}

type Tidal struct {
	API *tidalapi.TidalAPI
}

func (t *Tidal) GetAlbum(id string) (tonearm.Album, error) {
	album, err := t.API.OpenAPI.V2.Albums.Album(context.Background(), id, "artists.profileArt", "coverArt")
	if err != nil {
		return nil, err
	}

	return openapi.NewAlbum(*album), nil
}

func (t *Tidal) GetAlbumInfo(id string) (tonearm.AlbumInfo, error) {
	album, err := t.API.OpenAPI.V2.Albums.Album(context.Background(), id, "coverArt")
	if err != nil {
		return nil, err
	}

	return openapi.NewAlbumInfo(*album), nil
}

func (t *Tidal) GetAlbumTracks(id string) (tonearm.Paginator[tonearm.Track], error) {
	albumInfo, err := t.GetAlbumInfo(id)
	if err != nil {
		return nil, err
	}

	paginator := v1.NewPaginatorV1(t.API.V1.Albums.Items, id, func(pr *modelv1.PaginatedResponse[modelv1.AlbumItem]) []tonearm.Track {
		items := make([]tonearm.Track, len(pr.Items))
		for i, item := range pr.Items {
			items[i] = v1.NewTrack(item.Item, albumInfo)
		}
		return items
	})
	return paginator, nil
}

func (t *Tidal) GetArtist(id string) (tonearm.ArtistInfo, error) {
	artist, err := t.API.OpenAPI.V2.Artists.Artist(context.Background(), id, "profileArt")
	if err != nil {
		return nil, err
	}

	return openapi.NewArtistInfo(*artist), nil
}

func (t *Tidal) GetPlaylist(id string) (tonearm.Playlist, error) {
	playlist, err := t.API.OpenAPI.V2.Playlists.Playlist(context.Background(), id, "coverArt", "ownerProfiles.profileArt")
	if err != nil {
		return nil, err
	}

	return openapi.NewPlaylist(*playlist), nil
}

func (t *Tidal) GetPlaylistInfo(id string) (tonearm.PlaylistInfo, error) {
	playlist, err := t.API.OpenAPI.V2.Playlists.Playlist(context.Background(), id, "coverArt")
	if err != nil {
		return nil, err
	}

	return openapi.NewPlaylist(*playlist), nil
}

func (t *Tidal) GetPlaylistTracks(id string) (tonearm.Paginator[tonearm.Track], error) {
	paginator := openapi.NewPaginator(t.API.OpenAPI.V2.Playlists.Items, id, func(r *modelopenapi.Response[[]modelopenapi.Relationship]) []tonearm.Track {
		results := r.Included.Tracks(r.Data...)
		tracks := make([]tonearm.Track, len(results))
		for i, track := range results {
			tracks[i] = openapi.NewTrack(track)
		}
		return tracks
	}, "items.artists.profileArt", "items.albums.coverArt")
	return paginator, nil
}

func (t *Tidal) GetTrack(id string) (tonearm.Track, error) {
	track, err := t.API.OpenAPI.V2.Tracks.Track(context.Background(), id, "albums.coverArt", "artists.profileArt")
	if err != nil {
		return nil, err
	}

	return openapi.NewTrack(*track), nil
}

func NewTidal(api *tidalapi.TidalAPI) tonearm.Service {
	return &Tidal{API: api}
}
