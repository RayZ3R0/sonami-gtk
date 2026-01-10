package v1

import (
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/internal"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/v1/pages"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/v1/playlists"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/v1/tracks"
)

type V1 struct {
	client    *internal.Client
	Pages     *pages.Pages
	Playlists *playlists.Playlists
	Tracks    *tracks.Tracks
}

func New(client *internal.Client) *V1 {
	return &V1{
		client:    client,
		Pages:     pages.New(client),
		Playlists: playlists.New(client),
		Tracks:    tracks.New(client),
	}
}
