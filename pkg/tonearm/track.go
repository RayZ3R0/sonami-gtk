package tonearm

import "time"

type TrackQuality string

const (
	TrackQualityMax    TrackQuality = "max"
	TrackQualityHigh   TrackQuality = "high"
	TrackQualityMedium TrackQuality = "medium"
	TrackQualityLow    TrackQuality = "low"
)

// A track defines the minimum required information for a Track to be
// playable in Tonearm.
type Track interface {
	// Album returns (and possibly resolves) the album the track belongs to
	Album() (Album, error)

	// Artists returns (and possibly resolves) the artists featured on the track
	Artists() (Paginator[Artist], error)

	// ID returns the unique identifier for the track
	ID() string

	// Duration returns the total duration of the track
	Duration() time.Duration

	// Title returns the title of the track without any additional information
	Title() string

	// URL returns the shareable URL for the track
	URL() string
}
