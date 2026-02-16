package tonearm

type ArtistInfos []ArtistInfo

func (infos ArtistInfos) Names() []string {
	names := make([]string, len(infos))
	for i, info := range infos {
		names[i] = info.Name()
	}
	return names
}

type ArtistInfo interface {
	// ID returns the unique identifier for the artist
	ID() string

	// Name returns the name of the artist
	Name() string

	// ProfilePicture returns the URL of the artist's profile picture
	// If the artist has multiple profile pictures, the preferredSize parameter is used to determine which picture to return.
	ProfilePicture(preferredSize int) string

	// URL returns the shareable URL for the artist
	URL() string
}
