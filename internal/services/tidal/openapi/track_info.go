package openapi

import (
	"slices"
	"time"

	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/models/openapi"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
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

func (t TrackInfo) IsStreamable() bool {
	return slices.Contains(t.Data.Attributes.Availability, openapi.TrackAvailabilityStream)
}

func (t TrackInfo) Title() string {
	return t.Data.Attributes.Title
}

func (t TrackInfo) URL() string {
	return "https://tidal.com/track/" + t.ID()
}

func (t TrackInfo) Version() string {
	return t.Data.Attributes.Version
}

func NewTrackInfo(item openapi.Track) sonami.TrackInfo {
	return &TrackInfo{item}
}
