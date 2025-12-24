package tracks

import "codeberg.org/dergs/tidalwave/pkg/tidalapi/internal"

type Tracks struct {
	client *internal.Client
}

func New(client *internal.Client) *Tracks {
	return &Tracks{
		client: client,
	}
}
