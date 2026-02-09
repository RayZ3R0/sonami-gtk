package favourites

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"codeberg.org/dergs/tonearm/pkg/tidalapi/internal"
	v2 "codeberg.org/dergs/tonearm/pkg/tidalapi/models/v2"
)

type FavouriteArtist struct {
	client *internal.Client
}

func NewFavouriteArtist(client *internal.Client) *FavouriteArtist {
	return &FavouriteArtist{
		client: client,
	}
}

func (f *FavouriteArtist) Add(ctx context.Context, artistId string) error {
	body := url.Values{
		"artistIds":          {artistId},
		"onArtifactNotFound": {"FAIL"},
	}
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		"/v2/favourites/artists/add",
		strings.NewReader(body.Encode()),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		return err
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorText, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error while adding artist to favourites: %s", errorText)
	}

	return nil
}

func (f *FavouriteArtist) IDs(ctx context.Context, cursor, ifNoneMatch string) (data *v2.FavouritesIds, etag string, err error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"/v2/favourites/artists/ids",
		nil,
	)
	if err != nil {
		return nil, "", err
	}

	params := url.Values{}
	params.Add("limit", "50")
	if cursor != "" {
		params.Add("cursor", cursor)
	}
	req.URL.RawQuery = params.Encode()

	if ifNoneMatch != "" {
		req.Header.Set("If-None-Match", ifNoneMatch)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorText, _ := io.ReadAll(resp.Body)
		return nil, "", fmt.Errorf("error while retrieving favourite artist IDs: %s", errorText)
	}

	data = &v2.FavouritesIds{}
	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		return nil, "", err
	}

	return data, resp.Header.Get("ETag"), nil
}

func (f *FavouriteArtist) Remove(ctx context.Context, artistId string) error {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		"/v2/favourites/artists/remove",
		nil)
	if err != nil {
		return err
	}

	params := url.Values{}
	params.Add("artistIds", artistId)
	req.URL.RawQuery = params.Encode()

	resp, err := f.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err
	}

	return nil
}
