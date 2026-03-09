package state

import (
	"github.com/RayZ3R0/sonami-gtk/internal/localdb"
)

// FavouriteCache is the interface consumed by all UI components.
type FavouriteCache interface {
	Add(string) error
	Bust()
	Get() (*[]string, error)
	Remove(string) error
}

// Global caches backed by the local SQLite database.
var (
	AlbumsCache    FavouriteCache = localdb.NewFavouriteCache(localdb.FavouriteAlbum)
	ArtistsCache   FavouriteCache = localdb.NewFavouriteCache(localdb.FavouriteArtist)
	MixesCache     FavouriteCache = localdb.NewFavouriteCache(localdb.FavouriteMix)
	PlaylistsCache FavouriteCache = localdb.NewFavouriteCache(localdb.FavouritePlaylist)
	TracksCache    FavouriteCache = localdb.NewFavouriteCache(localdb.FavouriteTrack)
)
