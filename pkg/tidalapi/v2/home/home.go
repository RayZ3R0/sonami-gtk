package home

import (
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/internal"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/v2/home/feed"
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
