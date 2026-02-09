package favourites

import "codeberg.org/dergs/tonearm/pkg/tidalapi/internal"

type Favourites struct {
	client  *internal.Client
	Albums  *FavouriteAlbum
	Artists *FavouriteArtist
	Mixes   *FavouriteMixes
}

func New(client *internal.Client) *Favourites {
	return &Favourites{
		client:  client,
		Albums:  NewFavouriteAlbum(client),
		Artists: NewFavouriteArtist(client),
		Mixes:   NewFavouriteMixes(client),
	}
}
