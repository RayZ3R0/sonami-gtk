package tonearm

type Album interface {

	// Artists returns (and possibly resolves) the artists featured on the album
	Artists() (Paginator[Artist], error)

	// Cover returns (and possibly resolves) the cover image of the album
	//
	// If the album has multiple covers, the preferredSize parameter is used to determine which cover to return.
	Cover(preferredSize int) (string, error)

	// ID returns the unique identifier for the album
	ID() string

	// Title returns the title of the album
	Title() string

	// Tracks returns (and possibly resolves) the tracks on the album
	Tracks() (Paginator[Track], error)

	// URL returns the shareable URL for the album
	URL() string
}
