package v2

import (
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/internal"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/v2/home"
)

type V2 struct {
	client *internal.Client
	Home   *home.Home
}

func New(client *internal.Client) *V2 {
	return &V2{
		client: client,
		Home:   home.New(client),
	}
}
