package tidalapi

import (
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/internal"
	v1 "codeberg.org/dergs/tidalwave/pkg/tidalapi/v1"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/v2"
)

type TidalAPI struct {
	V1 *v1.V1
	V2 *v2.V2
}

func NewClient(countryCode string, token string) *TidalAPI {
	client := internal.NewClient(countryCode, token)

	return &TidalAPI{
		V1: v1.New(client),
		V2: v2.New(client),
	}
}
