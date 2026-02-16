package tonearm

import "time"

type AlbumInfo interface {
	PlaybackSource

	// Duration returns the duration of the album
	Duration() time.Duration

	// ID returns the unique identifier for the album
	ID() string

	// ReleasedAt returns the release date of the album
	ReleasedAt() time.Time

	// URL returns the shareable URL for the album
	URL() string
}

type Album interface {
	AlbumInfo

	Artists() ArtistInfos

	// Count returns the number of tracks in the album
	Count() int
}
