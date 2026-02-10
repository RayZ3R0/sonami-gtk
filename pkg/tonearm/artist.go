package tonearm

type Artist interface {
	// ID returns the unique identifier for the artist
	ID() string

	// Name returns the name of the artist
	Name() string

	// URL returns the shareable URL for the artist
	URL() string
}
