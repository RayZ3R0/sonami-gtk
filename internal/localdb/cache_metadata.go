package localdb

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

// CacheEntry represents a cached metadata entry
type CacheEntry struct {
	ID       string
	Data     []byte // JSON-serialized metadata
	CoverURL string
	CachedAt time.Time
	Error    sql.NullString // NULL if no error, error message if failed
}

// MetadataCache provides SQLite-backed caching for TIDAL metadata
type MetadataCache struct {
	tableName string
	logger    *slog.Logger
}

// NewAlbumCache creates cache for albums
func NewAlbumCache() *MetadataCache {
	return &MetadataCache{
		tableName: "cached_albums",
		logger:    slog.With("module", "cache", "type", "album"),
	}
}

// NewArtistCache creates cache for artists
func NewArtistCache() *MetadataCache {
	return &MetadataCache{
		tableName: "cached_artists",
		logger:    slog.With("module", "cache", "type", "artist"),
	}
}

// NewPlaylistCache creates cache for playlists
func NewPlaylistCache() *MetadataCache {
	return &MetadataCache{
		tableName: "cached_playlists",
		logger:    slog.With("module", "cache", "type", "playlist"),
	}
}

// NewTrackCache creates cache for tracks
func NewTrackCache() *MetadataCache {
	return &MetadataCache{
		tableName: "cached_tracks",
		logger:    slog.With("module", "cache", "type", "track"),
	}
}

// Get retrieves cached metadata by ID
func (c *MetadataCache) Get(id string) (*CacheEntry, error) {
	db := DB()

	var entry CacheEntry
	var cachedAtUnix int64

	query := fmt.Sprintf("SELECT id, data, cover_url, cached_at, error FROM %s WHERE id = ?", c.tableName)
	err := db.QueryRow(query, id).Scan(&entry.ID, &entry.Data, &entry.CoverURL, &cachedAtUnix, &entry.Error)

	if err == sql.ErrNoRows {
		return nil, nil // Not found (not an error)
	}
	if err != nil {
		return nil, err
	}

	entry.CachedAt = time.Unix(cachedAtUnix, 0)
	return &entry, nil
}

// GetBatch retrieves multiple entries efficiently using IN clause
func (c *MetadataCache) GetBatch(ids []string) (map[string]*CacheEntry, error) {
	if len(ids) == 0 {
		return make(map[string]*CacheEntry), nil
	}

	db := DB()
	result := make(map[string]*CacheEntry)

	// Build IN clause with placeholders
	query := fmt.Sprintf("SELECT id, data, cover_url, cached_at, error FROM %s WHERE id IN (", c.tableName)
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		if i > 0 {
			query += ", "
		}
		query += "?"
		args[i] = id
	}
	query += ")"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry CacheEntry
		var cachedAtUnix int64

		if err := rows.Scan(&entry.ID, &entry.Data, &entry.CoverURL, &cachedAtUnix, &entry.Error); err != nil {
			c.logger.Error("failed to scan cache entry", "error", err)
			continue
		}

		entry.CachedAt = time.Unix(cachedAtUnix, 0)
		result[entry.ID] = &entry
	}

	return result, rows.Err()
}

// Set stores metadata in cache
func (c *MetadataCache) Set(id string, data []byte, coverURL string) error {
	db := DB()

	query := fmt.Sprintf(`INSERT OR REPLACE INTO %s (id, data, cover_url, cached_at, error) 
	                      VALUES (?, ?, ?, ?, NULL)`, c.tableName)

	_, err := db.Exec(query, id, data, coverURL, time.Now().Unix())
	if err != nil {
		c.logger.Error("failed to cache metadata", "id", id, "error", err)
	}
	return err
}

// SetError marks an entry as failed (e.g., 404 from TIDAL)
func (c *MetadataCache) SetError(id string, errMsg string) error {
	db := DB()

	query := fmt.Sprintf(`INSERT OR REPLACE INTO %s (id, data, cover_url, cached_at, error) 
	                      VALUES (?, '', '', ?, ?)`, c.tableName)

	_, err := db.Exec(query, id, time.Now().Unix(), errMsg)
	if err != nil {
		c.logger.Error("failed to cache error", "id", id, "error", err)
	}
	return err
}

// Delete removes entry from cache
func (c *MetadataCache) Delete(id string) error {
	db := DB()

	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", c.tableName)
	_, err := db.Exec(query, id)
	return err
}

// GetStaleEntries returns IDs of entries cached more than `age` ago (for background sync)
// Only returns entries without errors (successful caches)
func (c *MetadataCache) GetStaleEntries(age time.Duration) ([]string, error) {
	db := DB()
	threshold := time.Now().Add(-age).Unix()

	query := fmt.Sprintf("SELECT id FROM %s WHERE cached_at < ? AND error IS NULL LIMIT 1000", c.tableName)
	rows, err := db.Query(query, threshold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// Clear removes all entries from this cache
func (c *MetadataCache) Clear() error {
	db := DB()

	query := fmt.Sprintf("DELETE FROM %s", c.tableName)
	_, err := db.Exec(query)
	if err != nil {
		c.logger.Error("failed to clear cache", "error", err)
	} else {
		c.logger.Info("cache cleared")
	}
	return err
}

// Count returns the number of cached entries
func (c *MetadataCache) Count() (int, error) {
	db := DB()

	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE error IS NULL", c.tableName)
	err := db.QueryRow(query).Scan(&count)
	return count, err
}
