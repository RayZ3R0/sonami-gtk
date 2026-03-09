package sonami

import "time"

type TrackQuality string

const (
	TrackQualityMax    TrackQuality = "max"
	TrackQualityHigh   TrackQuality = "high"
	TrackQualityMedium TrackQuality = "medium"
	TrackQualityLow    TrackQuality = "low"
)

type TrackInfo interface {
	Shareable

	// Duration returns the total duration of the track
	Duration() time.Duration

	// ID returns the unique identifier for the track
	ID() string

	// IsStreamable returns whether the track is streamable
	IsStreamable() bool

	// Title returns the title of the track without any additional information
	Title() string

	// Version returns the version of the track (e.g. "Acoustic")
	Version() string
}

type Track interface {
	TrackInfo
	PlaybackSource

	// Artists returns the artists associated with the track
	Artists() ArtistInfos

	// Album returns the album associated with the track
	Album() AlbumInfo
}

func FormatTitle(track TrackInfo) string {
	if track.Version() == "" {
		return track.Title()
	}
	return track.Title() + " (" + track.Version() + ")"
}
