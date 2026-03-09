package tidalapi

import (
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/auth"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/internal"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/openapi"
	v1 "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/v1"
	v2 "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/v2"
)

type TidalAPI struct {
	OpenAPI *openapi.OpenAPI
	V1      *v1.V1
	V2      *v2.V2
}

func NewClient(countryCode string, authStrategies ...auth.AuthStrategy) *TidalAPI {
	client := internal.NewClient(countryCode, authStrategies...)

	return &TidalAPI{
		OpenAPI: openapi.New(client),
		V1:      v1.New(client),
		V2:      v2.New(client),
	}
}
