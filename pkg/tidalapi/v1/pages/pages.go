package pages

import "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/internal"

type Pages struct {
	client *internal.Client
}

func New(client *internal.Client) *Pages {
	return &Pages{
		client: client,
	}
}
