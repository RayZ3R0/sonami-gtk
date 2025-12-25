package v2

import (
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/internal"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/openapi/v2/albums"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/openapi/v2/playlists"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/openapi/v2/tracks"
)

type V2 struct {
	Albums    *albums.Albums
	Playlists *playlists.Playlists
	Tracks    *tracks.Tracks
}

func New(client *internal.Client) *V2 {
	return &V2{
		Albums:    albums.New(client),
		Playlists: playlists.New(client),
		Tracks:    tracks.New(client),
	}
}
