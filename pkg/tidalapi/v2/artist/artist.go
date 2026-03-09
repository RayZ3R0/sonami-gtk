package artist

import (
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/internal"
)

type Artist struct {
	client *internal.Client
}

func New(client *internal.Client) *Artist {
	return &Artist{
		client: client,
	}
}
