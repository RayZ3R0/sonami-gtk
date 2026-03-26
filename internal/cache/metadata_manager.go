package cache

import (
	"encoding/json"
	"log/slog"

	"github.com/RayZ3R0/sonami-gtk/internal/localdb"
	"github.com/RayZ3R0/sonami-gtk/internal/services/tidal/openapi"
	"github.com/RayZ3R0/sonami-gtk/internal/services/tidal/v2"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
	modelopenapi "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/models/openapi"
	modelv2 "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/models/v2"
)

// MetadataManager coordinates L1 (memory) and L2 (SQLite) caches
// for album, artist, playlist, and track metadata
type MetadataManager struct {
	// L1 caches (memory, LRU with 500 capacity each)
	albumsL1    *LRUCache[string, sonami.Album]
	artistsL1   *LRUCache[string, sonami.Artist]
	playlistsL1 *LRUCache[string, sonami.Playlist]
	tracksL1    *LRUCache[string, sonami.Track]

	// L2 caches (SQLite)
	albumsL2    *localdb.MetadataCache
	artistsL2   *localdb.MetadataCache
	playlistsL2 *localdb.MetadataCache
	tracksL2    *localdb.MetadataCache

	logger *slog.Logger
}

// NewMetadataManager creates a new metadata cache manager
func NewMetadataManager() *MetadataManager {
	return &MetadataManager{
		albumsL1:    NewLRU[string, sonami.Album](500),
		artistsL1:   NewLRU[string, sonami.Artist](500),
		playlistsL1: NewLRU[string, sonami.Playlist](500),
		tracksL1:    NewLRU[string, sonami.Track](500),

		albumsL2:    localdb.NewAlbumCache(),
		artistsL2:   localdb.NewArtistCache(),
		playlistsL2: localdb.NewPlaylistCache(),
		tracksL2:    localdb.NewTrackCache(),

		logger: slog.With("module", "metadata-manager"),
	}
}

// GetAlbum retrieves album from cache (L1 → L2 → nil)
func (m *MetadataManager) GetAlbum(id string) (sonami.Album, bool) {
	// Try L1
	if album, ok := m.albumsL1.Get(id); ok {
		m.logger.Debug("cache hit L1", "type", "album", "id", id)
		return album, true
	}

	// Try L2
	entry, err := m.albumsL2.Get(id)
	if err != nil {
		m.logger.Error("L2 cache error", "type", "album", "id", id, "error", err)
		return nil, false
	}

	if entry != nil {
		// Check if cached error (404)
		if entry.Error.Valid {
			m.logger.Debug("cache hit L2 (error)", "type", "album", "id", id, "error", entry.Error.String)
			return nil, false
		}

		// Deserialize from JSON
		var albumData modelopenapi.Album
		if err := json.Unmarshal(entry.Data, &albumData); err != nil {
			m.logger.Warn("failed to deserialize album, clearing cache entry", "id", id, "error", err)
			// Delete corrupted cache entry and return miss
			_ = m.albumsL2.Delete(id)
			return nil, false
		}

		// Wrap in sonami interface
		album := openapi.NewAlbum(albumData)
		m.albumsL1.Set(id, album) // Populate L1
		m.logger.Debug("cache hit L2", "type", "album", "id", id)
		return album, true
	}

	m.logger.Debug("cache miss", "type", "album", "id", id)
	return nil, false
}

// SetAlbum stores album in both caches
func (m *MetadataManager) SetAlbum(id string, album sonami.Album, rawData modelopenapi.Album) error {
	// Serialize to JSON
	data, err := json.Marshal(rawData)
	if err != nil {
		m.logger.Error("failed to serialize album", "id", id, "error", err)
		return err
	}

	// Store in L2
	coverURL := album.Cover(172)
	if err := m.albumsL2.Set(id, data, coverURL); err != nil {
		m.logger.Error("failed to cache album in L2", "id", id, "error", err)
		// Continue anyway to populate L1
	}

	// Store in L1
	m.albumsL1.Set(id, album)
	m.logger.Debug("cached album", "id", id)
	return nil
}

// SetAlbumError marks album as failed (404)
func (m *MetadataManager) SetAlbumError(id string, errMsg string) error {
	return m.albumsL2.SetError(id, errMsg)
}

// InvalidateAlbum removes album from both caches
func (m *MetadataManager) InvalidateAlbum(id string) {
	m.albumsL1.Delete(id)
	m.albumsL2.Delete(id)
}

// GetArtist retrieves artist from cache (L1 → L2 → nil)
func (m *MetadataManager) GetArtist(id string) (sonami.Artist, bool) {
	// Try L1
	if artist, ok := m.artistsL1.Get(id); ok {
		m.logger.Debug("cache hit L1", "type", "artist", "id", id)
		return artist, true
	}

	// Try L2
	entry, err := m.artistsL2.Get(id)
	if err != nil {
		m.logger.Error("L2 cache error", "type", "artist", "id", id, "error", err)
		return nil, false
	}

	if entry != nil {
		if entry.Error.Valid {
			m.logger.Debug("cache hit L2 (error)", "type", "artist", "id", id, "error", entry.Error.String)
			return nil, false
		}

		var artistData modelv2.ArtistPage
		if err := json.Unmarshal(entry.Data, &artistData); err != nil {
			m.logger.Warn("failed to deserialize artist, clearing cache entry", "id", id, "error", err)
			_ = m.artistsL2.Delete(id)
			return nil, false
		}

		artist := v2.NewArtist(artistData)
		m.artistsL1.Set(id, artist)
		m.logger.Debug("cache hit L2", "type", "artist", "id", id)
		return artist, true
	}

	m.logger.Debug("cache miss", "type", "artist", "id", id)
	return nil, false
}

// SetArtist stores artist in both caches
func (m *MetadataManager) SetArtist(id string, artist sonami.Artist, rawData modelv2.ArtistPage) error {
	data, err := json.Marshal(rawData)
	if err != nil {
		m.logger.Error("failed to serialize artist", "id", id, "error", err)
		return err
	}

	coverURL := artist.Cover(172)
	if err := m.artistsL2.Set(id, data, coverURL); err != nil {
		m.logger.Error("failed to cache artist in L2", "id", id, "error", err)
	}

	m.artistsL1.Set(id, artist)
	m.logger.Debug("cached artist", "id", id)
	return nil
}

// SetArtistError marks artist as failed
func (m *MetadataManager) SetArtistError(id string, errMsg string) error {
	return m.artistsL2.SetError(id, errMsg)
}

// InvalidateArtist removes artist from both caches
func (m *MetadataManager) InvalidateArtist(id string) {
	m.artistsL1.Delete(id)
	m.artistsL2.Delete(id)
}

// GetPlaylist retrieves playlist from cache (L1 → L2 → nil)
func (m *MetadataManager) GetPlaylist(id string) (sonami.Playlist, bool) {
	// Try L1
	if playlist, ok := m.playlistsL1.Get(id); ok {
		m.logger.Debug("cache hit L1", "type", "playlist", "id", id)
		return playlist, true
	}

	// Try L2
	entry, err := m.playlistsL2.Get(id)
	if err != nil {
		m.logger.Error("L2 cache error", "type", "playlist", "id", id, "error", err)
		return nil, false
	}

	if entry != nil {
		if entry.Error.Valid {
			m.logger.Debug("cache hit L2 (error)", "type", "playlist", "id", id, "error", entry.Error.String)
			return nil, false
		}

		var playlistData modelopenapi.Playlist
		if err := json.Unmarshal(entry.Data, &playlistData); err != nil {
			m.logger.Warn("failed to deserialize playlist, clearing cache entry", "id", id, "error", err)
			_ = m.playlistsL2.Delete(id)
			return nil, false
		}

		playlist := openapi.NewPlaylist(playlistData)
		m.playlistsL1.Set(id, playlist)
		m.logger.Debug("cache hit L2", "type", "playlist", "id", id)
		return playlist, true
	}

	m.logger.Debug("cache miss", "type", "playlist", "id", id)
	return nil, false
}

// SetPlaylist stores playlist in both caches
func (m *MetadataManager) SetPlaylist(id string, playlist sonami.Playlist, rawData modelopenapi.Playlist) error {
	data, err := json.Marshal(rawData)
	if err != nil {
		m.logger.Error("failed to serialize playlist", "id", id, "error", err)
		return err
	}

	coverURL := playlist.Cover(172)
	if err := m.playlistsL2.Set(id, data, coverURL); err != nil {
		m.logger.Error("failed to cache playlist in L2", "id", id, "error", err)
	}

	m.playlistsL1.Set(id, playlist)
	m.logger.Debug("cached playlist", "id", id)
	return nil
}

// SetPlaylistError marks playlist as failed
func (m *MetadataManager) SetPlaylistError(id string, errMsg string) error {
	return m.playlistsL2.SetError(id, errMsg)
}

// InvalidatePlaylist removes playlist from both caches
func (m *MetadataManager) InvalidatePlaylist(id string) {
	m.playlistsL1.Delete(id)
	m.playlistsL2.Delete(id)
}

// GetTrack retrieves track from cache (L1 → L2 → nil)
func (m *MetadataManager) GetTrack(id string) (sonami.Track, bool) {
	// Try L1
	if track, ok := m.tracksL1.Get(id); ok {
		m.logger.Debug("cache hit L1", "type", "track", "id", id)
		return track, true
	}

	// Try L2
	entry, err := m.tracksL2.Get(id)
	if err != nil {
		m.logger.Error("L2 cache error", "type", "track", "id", id, "error", err)
		return nil, false
	}

	if entry != nil {
		if entry.Error.Valid {
			m.logger.Debug("cache hit L2 (error)", "type", "track", "id", id, "error", entry.Error.String)
			return nil, false
		}

		var trackData modelopenapi.Track
		if err := json.Unmarshal(entry.Data, &trackData); err != nil {
			m.logger.Warn("failed to deserialize track, clearing cache entry", "id", id, "error", err)
			_ = m.tracksL2.Delete(id)
			return nil, false
		}

		track := openapi.NewTrack(trackData)
		m.tracksL1.Set(id, track)
		m.logger.Debug("cache hit L2", "type", "track", "id", id)
		return track, true
	}

	m.logger.Debug("cache miss", "type", "track", "id", id)
	return nil, false
}

// SetTrack stores track in both caches
func (m *MetadataManager) SetTrack(id string, track sonami.Track, rawData modelopenapi.Track) error {
	data, err := json.Marshal(rawData)
	if err != nil {
		m.logger.Error("failed to serialize track", "id", id, "error", err)
		return err
	}

	coverURL := track.Cover(172)
	if err := m.tracksL2.Set(id, data, coverURL); err != nil {
		m.logger.Error("failed to cache track in L2", "id", id, "error", err)
	}

	m.tracksL1.Set(id, track)
	m.logger.Debug("cached track", "id", id)
	return nil
}

// SetTrackError marks track as failed
func (m *MetadataManager) SetTrackError(id string, errMsg string) error {
	return m.tracksL2.SetError(id, errMsg)
}

// InvalidateTrack removes track from both caches
func (m *MetadataManager) InvalidateTrack(id string) {
	m.tracksL1.Delete(id)
	m.tracksL2.Delete(id)
}

// ClearAll clears all caches
func (m *MetadataManager) ClearAll() error {
	m.albumsL1.Clear()
	m.artistsL1.Clear()
	m.playlistsL1.Clear()
	m.tracksL1.Clear()

	if err := m.albumsL2.Clear(); err != nil {
		return err
	}
	if err := m.artistsL2.Clear(); err != nil {
		return err
	}
	if err := m.playlistsL2.Clear(); err != nil {
		return err
	}
	if err := m.tracksL2.Clear(); err != nil {
		return err
	}

	m.logger.Info("all caches cleared")
	return nil
}

// Stats returns cache statistics
func (m *MetadataManager) Stats() CacheStats {
	albumsCount, _ := m.albumsL2.Count()
	artistsCount, _ := m.artistsL2.Count()
	playlistsCount, _ := m.playlistsL2.Count()
	tracksCount, _ := m.tracksL2.Count()

	return CacheStats{
		AlbumsL1:    m.albumsL1.Len(),
		AlbumsL2:    albumsCount,
		ArtistsL1:   m.artistsL1.Len(),
		ArtistsL2:   artistsCount,
		PlaylistsL1: m.playlistsL1.Len(),
		PlaylistsL2: playlistsCount,
		TracksL1:    m.tracksL1.Len(),
		TracksL2:    tracksCount,
	}
}

// CacheStats holds cache statistics
type CacheStats struct {
	AlbumsL1    int
	AlbumsL2    int
	ArtistsL1   int
	ArtistsL2   int
	PlaylistsL1 int
	PlaylistsL2 int
	TracksL1    int
	TracksL2    int
}
