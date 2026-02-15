package tidal

import (
	"context"
	"errors"
	"log/slog"

	"codeberg.org/dergs/tonearm/internal/secrets"
	"codeberg.org/dergs/tonearm/internal/services/tidal/openapi"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
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

func (t *Tidal) GetAlbum(id string, hints ...tonearm.FetchHint) (tonearm.Album, error) {
	included := openapi.AlbumHintsToInclude("", hints...)
	included = append(included, openapi.TrackHintsToInclude("items.", hints...)...)
	album, err := t.API.OpenAPI.V2.Albums.Album(context.Background(), id, included...)
	if err != nil {
		return nil, err
	}

	return openapi.NewAlbum(t.API, *album), nil
}

func (t *Tidal) GetArtist(id string, hints ...tonearm.FetchHint) (tonearm.Artist, error) {
	return nil, errors.New("not implemented")
}

func (t *Tidal) GetTrack(id string, hints ...tonearm.FetchHint) (tonearm.Track, error) {
	return nil, errors.New("not implemented")
}

func NewTidal(api *tidalapi.TidalAPI) tonearm.Service {
	return &Tidal{API: api}
}
