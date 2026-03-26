package localdb

import (
	"database/sql"
	"time"
)

// PageFeedEntry represents a cached page feed
type PageFeedEntry struct {
	ID       string
	Data     []byte
	CachedAt time.Time
}

// GetPageFeed retrieves a cached page feed by ID
func GetPageFeed(id string) (*PageFeedEntry, error) {
	var entry PageFeedEntry
	var cachedAtUnix int64

	err := DB().QueryRow(
		`SELECT id, data, cached_at FROM cached_page_feeds WHERE id = ?`,
		id,
	).Scan(&entry.ID, &entry.Data, &cachedAtUnix)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	entry.CachedAt = time.Unix(cachedAtUnix, 0)
	return &entry, nil
}

// SetPageFeed stores a page feed in the cache
func SetPageFeed(id string, data []byte) error {
	_, err := DB().Exec(
		`INSERT OR REPLACE INTO cached_page_feeds (id, data, cached_at) VALUES (?, ?, ?)`,
		id, data, time.Now().Unix(),
	)
	return err
}

// DeletePageFeed removes a page feed from the cache
func DeletePageFeed(id string) error {
	_, err := DB().Exec(`DELETE FROM cached_page_feeds WHERE id = ?`, id)
	return err
}

// ClearAllPageFeeds removes all cached page feeds
func ClearAllPageFeeds() error {
	_, err := DB().Exec(`DELETE FROM cached_page_feeds`)
	return err
}
