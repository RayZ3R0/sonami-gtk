package albums

import "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/internal"

type Albums struct {
	client *internal.Client
}

func New(client *internal.Client) *Albums {
	return &Albums{
		client: client,
	}
}
