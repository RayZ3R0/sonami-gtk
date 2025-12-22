package home

import (
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/internal"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/v2/home/feed"
)

type Home struct {
	client *internal.Client
	Feed   *feed.Feed
}

func New(client *internal.Client) *Home {
	return &Home{
		client: client,
		Feed:   feed.New(client),
	}
}
