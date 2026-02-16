package tonearm

type ArtistInfos []ArtistInfo

func (infos ArtistInfos) Names() []string {
	names := make([]string, len(infos))
	for i, info := range infos {
		names[i] = info.Title()
	}
	return names
}

type ArtistInfo interface {
	PlaybackSource

	// ID returns the unique identifier for the artist
	ID() string

	// URL returns the shareable URL for the artist
	URL() string
}
