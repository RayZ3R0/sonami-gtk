package tracks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"codeberg.org/dergs/tidalwave/pkg/tidalapi/helper"
	v1 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v1"
)

type PlaybackInfoOptions struct {
	AssetPresentation v1.AssetPresentation `json:"assetPresentation"`
	AudioQuality      v1.AudioQuality      `json:"audioQuality"`
	PlaybackMode      v1.PlaybackMode      `json:"playbackMode"`
}

func (p *Tracks) PlaybackInfo(ctx context.Context, trackId string, opts PlaybackInfoOptions) (*v1.PlaybackInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("/v1/tracks/%s/playbackinfo", trackId), nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Set("assetpresentation", helper.OptionalString(string(opts.AssetPresentation), string(v1.AssetPresentationFull)))
	params.Set("audioquality", helper.OptionalString(string(opts.AudioQuality), string(v1.AudioQualityHighResLossess)))
	params.Set("playbackmode", helper.OptionalString(string(opts.PlaybackMode), string(v1.PlaybackModeStream)))
	req.URL.RawQuery = params.Encode()

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var playbackInfo v1.PlaybackInfo
	if err := json.NewDecoder(resp.Body).Decode(&playbackInfo); err != nil {
		return nil, err
	}

	return &playbackInfo, nil
}
