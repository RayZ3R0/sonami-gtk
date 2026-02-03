package feed

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"codeberg.org/dergs/tonearm/pkg/tidalapi/internal"
)

type Feed struct {
	client *internal.Client
}

func New(client *internal.Client) *Feed {
	return &Feed{
		client: client,
	}
}

func (f *Feed) Activities(ctx context.Context, id string) ([]*Activity, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/v2/feed/activities", nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Set("userId", id)
	req.URL.RawQuery = params.Encode()

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var activities []*Activity
	response := struct {
		Activities *[]*Activity `json:"activities"`
	}{
		Activities: &activities,
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return activities, nil
}
