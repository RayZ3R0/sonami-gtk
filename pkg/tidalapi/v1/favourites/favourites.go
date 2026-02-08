package favourites

import "codeberg.org/dergs/tonearm/pkg/tidalapi/internal"

type Favourites struct {
	client *internal.Client
}

func New(client *internal.Client) *Favourites {
	return &Favourites{
		client: client,
	}
}
