package cache

import (
	"context"
	"log/slog"
	"sync/atomic"
	"time"
)

// SyncManager handles background sync of stale cache entries (>7 days old)
type SyncManager struct {
	service        *CachedService
	syncInProgress atomic.Bool
	logger         *slog.Logger
}

// NewSyncManager creates a new sync manager
func NewSyncManager(service *CachedService) *SyncManager {
	return &SyncManager{
		service: service,
		logger:  slog.With("module", "sync-manager"),
	}
}

// SyncStaleEntries checks for and updates stale cache entries (>7 days old)
// This runs in the background and rate-limits to 1 request per 100ms
func (m *SyncManager) SyncStaleEntries(ctx context.Context) {
	// Only allow one sync at a time
	if !m.syncInProgress.CompareAndSwap(false, true) {
		m.logger.Debug("sync already in progress, skipping")
		return
	}
	defer m.syncInProgress.Store(false)

	m.logger.Info("starting background sync")
	startTime := time.Now()

	// Get stale entries (cached > 7 days ago)
	staleAge := 7 * 24 * time.Hour

	totalSynced := 0

	// Sync albums
	if staleAlbums, err := m.service.cache.albumsL2.GetStaleEntries(staleAge); err == nil && len(staleAlbums) > 0 {
		m.logger.Info("syncing stale albums", "count", len(staleAlbums))
		synced := m.syncAlbums(ctx, staleAlbums)
		totalSynced += synced
		m.logger.Debug("albums synced", "count", synced)
	}

	// Sync artists
	if staleArtists, err := m.service.cache.artistsL2.GetStaleEntries(staleAge); err == nil && len(staleArtists) > 0 {
		m.logger.Info("syncing stale artists", "count", len(staleArtists))
		synced := m.syncArtists(ctx, staleArtists)
		totalSynced += synced
		m.logger.Debug("artists synced", "count", synced)
	}

	// Sync playlists
	if stalePlaylists, err := m.service.cache.playlistsL2.GetStaleEntries(staleAge); err == nil && len(stalePlaylists) > 0 {
		m.logger.Info("syncing stale playlists", "count", len(stalePlaylists))
		synced := m.syncPlaylists(ctx, stalePlaylists)
		totalSynced += synced
		m.logger.Debug("playlists synced", "count", synced)
	}

	// Sync tracks
	if staleTracks, err := m.service.cache.tracksL2.GetStaleEntries(staleAge); err == nil && len(staleTracks) > 0 {
		m.logger.Info("syncing stale tracks", "count", len(staleTracks))
		synced := m.syncTracks(ctx, staleTracks)
		totalSynced += synced
		m.logger.Debug("tracks synced", "count", synced)
	}

	duration := time.Since(startTime)
	m.logger.Info("background sync completed", "total_synced", totalSynced, "duration_ms", duration.Milliseconds())
}

func (m *SyncManager) syncAlbums(ctx context.Context, ids []string) int {
	ticker := time.NewTicker(100 * time.Millisecond) // Rate limit: 1 req / 100ms
	defer ticker.Stop()

	synced := 0
	for _, id := range ids {
		select {
		case <-ctx.Done():
			return synced // Cancelled
		case <-ticker.C:
			// Invalidate L1 cache so it re-fetches
			m.service.cache.albumsL1.Delete(id)

			// Re-fetch from TIDAL (will update cache)
			if _, err := m.service.GetAlbum(id); err != nil {
				m.logger.Debug("failed to sync album", "id", id, "error", err)
			} else {
				synced++
			}
		}
	}
	return synced
}

func (m *SyncManager) syncArtists(ctx context.Context, ids []string) int {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	synced := 0
	for _, id := range ids {
		select {
		case <-ctx.Done():
			return synced
		case <-ticker.C:
			m.service.cache.artistsL1.Delete(id)

			if _, err := m.service.GetArtist(id); err != nil {
				m.logger.Debug("failed to sync artist", "id", id, "error", err)
			} else {
				synced++
			}
		}
	}
	return synced
}

func (m *SyncManager) syncPlaylists(ctx context.Context, ids []string) int {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	synced := 0
	for _, id := range ids {
		select {
		case <-ctx.Done():
			return synced
		case <-ticker.C:
			m.service.cache.playlistsL1.Delete(id)

			if _, err := m.service.GetPlaylist(id); err != nil {
				m.logger.Debug("failed to sync playlist", "id", id, "error", err)
			} else {
				synced++
			}
		}
	}
	return synced
}

func (m *SyncManager) syncTracks(ctx context.Context, ids []string) int {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	synced := 0
	for _, id := range ids {
		select {
		case <-ctx.Done():
			return synced
		case <-ticker.C:
			m.service.cache.tracksL1.Delete(id)

			if _, err := m.service.GetTrack(id); err != nil {
				m.logger.Debug("failed to sync track", "id", id, "error", err)
			} else {
				synced++
			}
		}
	}
	return synced
}

// IsRunning returns true if a sync is currently in progress
func (m *SyncManager) IsRunning() bool {
	return m.syncInProgress.Load()
}
