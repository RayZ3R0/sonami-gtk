package services

import (
	"github.com/RayZ3R0/sonami-gtk/internal/services/tidal"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
	"github.com/infinytum/injector"
)

func init() {
	injector.DeferredSingleton(func(api *tidalapi.TidalAPI) sonami.Service {
		return tidal.NewTidal(api)
	})
}
