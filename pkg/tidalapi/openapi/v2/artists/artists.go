package artists

import "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/internal"

type Artists struct {
	client *internal.Client
}

func New(client *internal.Client) *Artists {
	return &Artists{
		client: client,
	}
}
