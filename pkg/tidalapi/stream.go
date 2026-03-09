package tidalapi

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"

	v1 "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/models/v1"
)

var streamLogger = slog.With("module", "stream_resolver")

// streamAPIResponse is the top-level wrapper from the new API.
// All actual data is under the "data" key.
type streamAPIResponse struct {
	Version string         `json:"version"`
	Data    streamResponse `json:"data"`
}

// streamResponse contains the actual stream data returned by the new API.
// Matches the /track/?id=X&quality=Y response.
type streamResponse struct {
	TrackID           int    `json:"trackId"`
	AssetPresentation string `json:"assetPresentation"`
	AudioMode         string `json:"audioMode"`
	AudioQuality      string `json:"audioQuality"`
	ManifestMimeType  string `json:"manifestMimeType"`
	Manifest          string `json:"manifest"`

	// ReplayGain — provided by the API (same fields as Tidal)
	AlbumReplayGain    *float64 `json:"albumReplayGain"`
	AlbumPeakAmplitude *float64 `json:"albumPeakAmplitude"`
	TrackReplayGain    *float64 `json:"trackReplayGain"`
	TrackPeakAmplitude *float64 `json:"trackPeakAmplitude"`

	// Direct URL fallback fields (defensive)
	URL         string `json:"url"`
	StreamURL   string `json:"streamUrl"`
	PlaybackURL string `json:"playbackUrl"`
}

// StreamResolver resolves track stream URLs via the new account-free API endpoints.
type StreamResolver struct {
	endpoints *EndpointManager
}

// NewStreamResolver creates a new StreamResolver backed by the given EndpointManager.
func NewStreamResolver(em *EndpointManager) *StreamResolver {
	return &StreamResolver{endpoints: em}
}

// Resolve fetches stream information for a track from the new API.
// It returns a v1.PlaybackInfo-compatible struct to minimize changes in the player layer.
func (sr *StreamResolver) Resolve(ctx context.Context, trackID string, quality v1.AudioQuality) (*v1.PlaybackInfo, error) {
	// Pass the quality directly to the API — it handles fallback server-side,
	// returning whatever quality the track supports (same as original Tidal behavior).
	path := fmt.Sprintf("/track/?id=%s&quality=%s", trackID, string(quality))

	resp, err := sr.endpoints.DoRequest(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("stream resolution failed for track %s: %w", trackID, err)
	}
	defer resp.Body.Close()

	// Decode the wrapped response: { "version": "...", "data": { ... } }
	var apiResp streamAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode stream response for track %s: %w", trackID, err)
	}

	return sr.toPlaybackInfo(apiResp.Data)
}

// toPlaybackInfo converts the stream response into a v1.PlaybackInfo struct.
// This allows the existing player pipeline (enqueue → playbin) to work without changes.
func (sr *StreamResolver) toPlaybackInfo(resp streamResponse) (*v1.PlaybackInfo, error) {
	// Use the quality and mode returned by the API directly
	info := &v1.PlaybackInfo{
		AudioQuality:      v1.AudioQuality(resp.AudioQuality),
		AudioMode:         v1.AudioMode(resp.AudioMode),
		AssetPresentation: v1.AssetPresentation(resp.AssetPresentation),

		// ReplayGain fields — populated from the API response
		AlbumReplayGain:    resp.AlbumReplayGain,
		AlbumPeakAmplitude: resp.AlbumPeakAmplitude,
		TrackReplayGain:    resp.TrackReplayGain,
		TrackPeakAmplitude: resp.TrackPeakAmplitude,
	}

	// If a manifest is present, use it
	if resp.Manifest != "" {
		info.Manifest = resp.Manifest
		if resp.ManifestMimeType != "" {
			info.ManifestMimeType = v1.ManifestMimeType(resp.ManifestMimeType)
		} else {
			info.ManifestMimeType = v1.ManifestMimeTypeAudioBTS
		}

		// Only validate BTS manifests (base64-encoded JSON with urls array).
		// MPD/DASH manifests are XML and handled by enqueueMPDStream() in the player layer.
		if info.ManifestMimeType == v1.ManifestMimeTypeAudioBTS {
			decoded, err := base64.StdEncoding.DecodeString(resp.Manifest)
			if err != nil {
				return nil, fmt.Errorf("failed to decode BTS manifest: %w", err)
			}

			var bts struct {
				URLs []string `json:"urls"`
			}
			if err := json.Unmarshal(decoded, &bts); err != nil {
				return nil, fmt.Errorf("failed to parse BTS manifest JSON: %w", err)
			}

			if len(bts.URLs) == 0 {
				return nil, fmt.Errorf("BTS manifest contains no URLs")
			}

			streamLogger.Debug("resolved stream via BTS manifest", "url_count", len(bts.URLs),
				"replay_gain_track", resp.TrackReplayGain)
		} else {
			streamLogger.Debug("resolved stream via MPD/DASH manifest",
				"quality", resp.AudioQuality, "replay_gain_track", resp.TrackReplayGain)
		}

		return info, nil
	}

	// Try direct URL fields (defensive fallback)
	directURL := resp.URL
	if directURL == "" {
		directURL = resp.StreamURL
	}
	if directURL == "" {
		directURL = resp.PlaybackURL
	}

	if directURL != "" {
		// Wrap the direct URL in a BTS-compatible manifest so the existing
		// enqueueBTSStream() function can handle it without modification
		btsJSON, _ := json.Marshal(map[string]any{
			"urls": []string{directURL},
		})
		info.Manifest = base64.StdEncoding.EncodeToString(btsJSON)
		info.ManifestMimeType = v1.ManifestMimeTypeAudioBTS

		streamLogger.Debug("resolved stream via direct URL")
		return info, nil
	}

	return nil, fmt.Errorf("stream response contains neither manifest nor direct URL")
}


