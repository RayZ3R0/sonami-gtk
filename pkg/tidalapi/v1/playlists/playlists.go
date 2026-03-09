package playlists

import "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/internal"

type Playlists struct {
	client *internal.Client
}

func New(client *internal.Client) *Playlists {
	return &Playlists{
		client: client,
	}
}
