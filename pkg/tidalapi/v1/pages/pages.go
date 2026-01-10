package pages

import "codeberg.org/dergs/tidalwave/pkg/tidalapi/internal"

type Pages struct {
	client *internal.Client
}

func New(client *internal.Client) *Pages {
	return &Pages{
		client: client,
	}
}
