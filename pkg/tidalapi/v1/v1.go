package v1

import (
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/internal"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/v1/playlists"
)

type V1 struct {
	client    *internal.Client
	Playlists *playlists.Playlists
}

func New(client *internal.Client) *V1 {
	return &V1{
		client:    client,
		Playlists: playlists.New(client),
	}
}
