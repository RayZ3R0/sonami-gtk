package tracks

import "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/internal"

type Tracks struct {
	client *internal.Client
}

func New(client *internal.Client) *Tracks {
	return &Tracks{
		client: client,
	}
}
