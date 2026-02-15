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
	Album(hints ...FetchHint) (Album, error)

	// ArtistNames returns the names of the artists featured on the track
	// This function is separate to Artists in case the backing service can
	// perform an optimized lookup for just the names.
	ArtistNames() ([]string, error)

	// Artists returns (and possibly resolves) the artists featured on the track
	Artists(hints ...FetchHint) (Paginator[Artist], error)

	// CoverURL returns the URL of the track's cover art
	// This convenience function is provided for easy access to the cover art URL.
	// If the track has multiple covers, the preferredSize parameter is used to determine which cover to return.
	Cover(preferredSize int) (string, error)

	// ID returns the unique identifier for the track
	ID() string

	// Duration returns the total duration of the track
	Duration() time.Duration

	// Title returns the title of the track without any additional information
	Title() string

	// URL returns the shareable URL for the track
	URL() string
}
