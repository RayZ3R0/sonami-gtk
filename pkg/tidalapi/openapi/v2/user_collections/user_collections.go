package user_collections

import "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/internal"

type UserCollections struct {
	client *internal.Client
}

func New(client *internal.Client) *UserCollections {
	return &UserCollections{
		client: client,
	}
}
