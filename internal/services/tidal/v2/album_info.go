package v2

import (
	"strconv"
	"time"

	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi"
	v2 "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/models/v2"
)

var albumInfoLogger = logger.With("type", "AlbumInfo").WithGroup("album_info")

type AlbumInfo struct {
	v2.AlbumItemData
}

func (a AlbumInfo) Cover(preferredSize int) string {
	logger := albumInfoLogger.With("method", "Cover").WithGroup("cover").With("preferred_size", preferredSize)

	if preferredSize < 0 {
		logger.Debug("legacy api does not support preferred size")
	}

	return tidalapi.ImageURL(a.AlbumItemData.Cover)
}

func (a AlbumInfo) Duration() time.Duration {
	return time.Duration(a.AlbumItemData.Duration) * time.Second
}

func (a AlbumInfo) ID() string {
	return strconv.Itoa(a.AlbumItemData.Id)
}

func (a AlbumInfo) ReleasedAt() time.Time {
	return a.AlbumItemData.ReleaseDate.Time
}

func (a AlbumInfo) Route() string {
	return "album/" + a.ID()
}

func (a AlbumInfo) SourceType() sonami.SourceType {
	return sonami.SourceTypeAlbum
}

func (a AlbumInfo) Title() string {
	return a.AlbumItemData.Title
}

func (a AlbumInfo) URL() string {
	return "https://tidal.com/album/" + a.ID()
}

func NewAlbumInfo(album v2.AlbumItemData) sonami.AlbumInfo {
	return &AlbumInfo{album}
}
