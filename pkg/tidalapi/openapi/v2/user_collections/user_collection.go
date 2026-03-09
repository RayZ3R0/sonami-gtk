package user_collections

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/models/openapi"
)

func (p *UserCollections) UserCollection(ctx context.Context, id string, include ...string) (*openapi.UserCollection, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://openapi.tidal.com/v2/userCollections/%s", id), nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Set("include", strings.Join(include, ","))
	req.URL.RawQuery = params.Encode()

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var userCollection openapi.UserCollection
	if err := json.NewDecoder(resp.Body).Decode(&userCollection); err != nil {
		return nil, err
	}

	return &userCollection, nil
}
