package tonearm

import "time"

type PlaylistInfo interface {
	PlaybackSource

	// CreatedAt returns the creation date of the playlist
	CreatedAt() time.Time

	// Duration returns the duration of the album
	Duration() time.Duration

	// ID returns the unique identifier for the album
	ID() string

	// IsMix returns true if the playlist is a mix
	IsMix() bool

	// URL returns the shareable URL for the album
	URL() string
}

type Playlist interface {
	PlaylistInfo

	// Count returns the number of tracks in the playlist
	Count() int

	Creator() ArtistInfo
}
