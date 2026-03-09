package openapi

import (
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/internal"
	v2 "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/openapi/v2"
)

type OpenAPI struct {
	V2 *v2.V2
}

func New(client *internal.Client) *OpenAPI {
	return &OpenAPI{
		V2: v2.New(client),
	}
}
