package openapi

import (
	"time"

	"codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
)

type TrackInfo struct {
	openapi.Track
}

func (t TrackInfo) Duration() time.Duration {
	return t.Data.Attributes.Duration.Duration
}

func (t TrackInfo) ID() string {
	return t.Data.ID
}

func (t TrackInfo) Title() string {
	return t.Data.Attributes.Title
}

func (t TrackInfo) URL() string {
	return "https://tidal.com/track/" + t.ID()
}

func NewTrackInfo(item openapi.Track) tonearm.TrackInfo {
	return &TrackInfo{item}
}
