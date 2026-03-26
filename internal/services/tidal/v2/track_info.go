package v2

import (
	"strconv"
	"time"

	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
	v2 "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/models/v2"
)

type TrackInfo struct {
	v2.TrackItemData
}

func (t TrackInfo) Duration() time.Duration {
	return time.Duration(t.TrackItemData.Duration) * time.Second
}

func (t TrackInfo) ID() string {
	return strconv.Itoa(t.TrackItemData.ID)
}

func (t TrackInfo) IsStreamable() bool {
	return t.AllowStreaming
}

func (t TrackInfo) Title() string {
	return t.TrackItemData.Title
}

func (t TrackInfo) URL() string {
	return "https://tidal.com/track/" + t.ID()
}

func (t TrackInfo) Version() string {
	return t.TrackItemData.Version
}

func NewTrackInfo(item v2.TrackItemData) sonami.TrackInfo {
	return &TrackInfo{TrackItemData: item}
}
