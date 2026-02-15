package tonearm

type FetchHint int

const (
	AlbumHintArtists FetchHint = iota
	AlbumHintCover
	AlbumHintTracks

	TrackHintAlbum FetchHint = iota
	TrackHintArtists

	ArtistHintAlbums FetchHint = iota
	ArtistHintProfilePicture
)

type Service interface {
	GetAlbum(id string, hints ...FetchHint) (Album, error)
	GetArtist(id string, hints ...FetchHint) (Artist, error)
	GetTrack(id string, hints ...FetchHint) (Track, error)
}
