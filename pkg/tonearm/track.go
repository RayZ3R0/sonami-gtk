package tonearm

import "time"

type TrackQuality string

const (
	TrackQualityMax    TrackQuality = "max"
	TrackQualityHigh   TrackQuality = "high"
	TrackQualityMedium TrackQuality = "medium"
	TrackQualityLow    TrackQuality = "low"
)

type TrackInfo interface {
	// Duration returns the total duration of the track
	Duration() time.Duration

	// ID returns the unique identifier for the track
	ID() string

	// IsStreamable returns whether the track is streamable
	IsStreamable() bool

	// Title returns the title of the track without any additional information
	Title() string

	// URL returns the shareable URL for the track
	URL() string
}

type Track interface {
	TrackInfo

	// Artists returns the artists associated with the track
	Artists() ArtistInfos

	// Album returns the album associated with the track
	Album() AlbumInfo

	// Cover returns the URL of the album cover for the track
	// If the album has multiple covers, the preferredSize parameter is used to determine which cover to return.
	Cover(preferredSize int) string
}
