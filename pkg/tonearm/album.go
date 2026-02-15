package tonearm

type Album interface {
	PlaybackSource

	// Artists returns (and possibly resolves) the artists featured on the album
	Artists(hints ...FetchHint) (Paginator[Artist], error)

	// ID returns the unique identifier for the album
	ID() string

	// Tracks returns (and possibly resolves) the tracks on the album
	Tracks(hints ...FetchHint) (Paginator[Track], error)

	// URL returns the shareable URL for the album
	URL() string
}
