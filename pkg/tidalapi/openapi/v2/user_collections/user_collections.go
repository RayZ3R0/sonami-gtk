package user_collections

import "codeberg.org/dergs/tidalwave/pkg/tidalapi/internal"

type UserCollections struct {
	client *internal.Client
}

func New(client *internal.Client) *UserCollections {
	return &UserCollections{
		client: client,
	}
}
