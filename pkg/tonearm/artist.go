package tonearm

type Artist interface {
	// ID returns the unique identifier for the artist
	ID() string

	// Name returns the name of the artist
	Name() string

	// ProfilePicture returns the URL of the artist's profile picture
	//
	// If the artist has multiple profile pictures, the preferredSize parameter is used to determine which picture to return.
	ProfilePicture(preferredSize int) (string, error)

	// URL returns the shareable URL for the artist
	URL() string
}
