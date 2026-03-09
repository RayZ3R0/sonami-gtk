package v1

import (
	"strconv"
	"time"

	v1 "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/models/v1"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
)

type TrackInfo struct {
	v1.AlbumItem
}

func (t TrackInfo) Duration() time.Duration {
	return time.Duration(t.AlbumItem.Duration) * time.Second
}

func (t TrackInfo) ID() string {
	return strconv.Itoa(t.AlbumItem.ID)
}

func (t TrackInfo) IsStreamable() bool {
	return t.AllowStreaming
}

func (t TrackInfo) Title() string {
	return t.AlbumItem.Title
}

func (t TrackInfo) URL() string {
	return "https://tidal.com/track/" + t.ID()
}

func (t TrackInfo) Version() string {
	return t.AlbumItem.Version
}

func NewTrackInfo(item v1.AlbumItem) sonami.TrackInfo {
	return &TrackInfo{AlbumItem: item}
}
