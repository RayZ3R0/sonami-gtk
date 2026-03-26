package services

import (
	"github.com/RayZ3R0/sonami-gtk/internal/cache"
	"github.com/RayZ3R0/sonami-gtk/internal/services/tidal"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi"
	"github.com/RayZ3R0/sonami-gtk/pkg/utils/imgutil"
	"github.com/infinytum/injector"
)

func init() {
	// Register the underlying TIDAL service (not directly exposed)
	injector.DeferredSingleton(func(api *tidalapi.TidalAPI) *tidal.Tidal {
		return tidal.NewTidal(api).(*tidal.Tidal)
	})

	// Register the cached service wrapper (implements sonami.Service)
	injector.DeferredSingleton(func(tidalService *tidal.Tidal) sonami.Service {
		return cache.NewCachedService(tidalService)
	})

	// Register the CachedService directly for sync manager access
	injector.DeferredSingleton(func(service sonami.Service) *cache.CachedService {
		cs := service.(*cache.CachedService)
		// Register hook to prefetch newly favourited items
		cs.RegisterFavouriteHook()
		return cs
	})

	// Register sync manager
	injector.DeferredSingleton(func(cachedService *cache.CachedService) *cache.SyncManager {
		return cache.NewSyncManager(cachedService)
	})

	// Register image prefetcher and wire it back to cached service
	injector.DeferredSingleton(func(imgUtil *imgutil.ImgUtil, cachedService *cache.CachedService) *cache.ImagePrefetcher {
		prefetcher := cache.NewImagePrefetcher(imgUtil)
		prefetcher.Start() // Start background worker immediately

		// Wire prefetcher into cached service
		cachedService.SetImagePrefetcher(prefetcher)

		return prefetcher
	})
}
