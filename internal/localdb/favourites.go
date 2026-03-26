package localdb

import (
	"log/slog"
	"sync"
)

// FavouriteType represents the category of a favourited resource.
type FavouriteType string

const (
	FavouriteAlbum    FavouriteType = "album"
	FavouriteArtist   FavouriteType = "artist"
	FavouriteMix      FavouriteType = "mix"
	FavouritePlaylist FavouriteType = "playlist"
	FavouriteTrack    FavouriteType = "track"
)

// FavouriteAddHook is called when a favourite is added.
// The hook receives the favourite type and ID.
type FavouriteAddHook func(favType FavouriteType, id string)

// favouriteAddHooks stores registered hooks for favourite additions
var (
	favouriteAddHooks   []FavouriteAddHook
	favouriteAddHooksMu sync.RWMutex
)

// RegisterFavouriteAddHook registers a hook to be called when a favourite is added.
// This is used by the cache system to prefetch newly favourited items.
func RegisterFavouriteAddHook(hook FavouriteAddHook) {
	favouriteAddHooksMu.Lock()
	defer favouriteAddHooksMu.Unlock()
	favouriteAddHooks = append(favouriteAddHooks, hook)
}

// notifyFavouriteAdded calls all registered hooks when a favourite is added
func notifyFavouriteAdded(favType FavouriteType, id string) {
	favouriteAddHooksMu.RLock()
	hooks := favouriteAddHooks
	favouriteAddHooksMu.RUnlock()

	for _, hook := range hooks {
		go hook(favType, id)
	}
}

// LocalFavouriteCache implements state.FavouriteCache backed by SQLite.
type LocalFavouriteCache struct {
	favType FavouriteType
	mu      sync.RWMutex
	cached  *[]string
}

// NewFavouriteCache creates a new local favourite cache for the given type.
func NewFavouriteCache(t FavouriteType) *LocalFavouriteCache {
	return &LocalFavouriteCache{favType: t}
}

func (c *LocalFavouriteCache) Add(id string) error {
	db := DB()

	_, err := db.Exec(
		`INSERT OR IGNORE INTO favourites (type, id) VALUES (?, ?)`,
		string(c.favType), id,
	)
	if err != nil {
		slog.Error("failed to add favourite", "type", c.favType, "id", id, "error", err)
		return err
	}

	c.Bust()

	// Notify hooks (e.g., cache prefetcher) about the new favourite
	notifyFavouriteAdded(c.favType, id)

	return nil
}

func (c *LocalFavouriteCache) Remove(id string) error {
	db := DB()

	_, err := db.Exec(
		`DELETE FROM favourites WHERE type = ? AND id = ?`,
		string(c.favType), id,
	)
	if err != nil {
		slog.Error("failed to remove favourite", "type", c.favType, "id", id, "error", err)
		return err
	}

	c.Bust()
	return nil
}

func (c *LocalFavouriteCache) Get() (*[]string, error) {
	c.mu.RLock()
	if c.cached != nil {
		defer c.mu.RUnlock()
		return c.cached, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check after acquiring write lock.
	if c.cached != nil {
		return c.cached, nil
	}

	db := DB()
	rows, err := db.Query(
		`SELECT id FROM favourites WHERE type = ? ORDER BY added_at DESC`,
		string(c.favType),
	)
	if err != nil {
		slog.Error("failed to query favourites", "type", c.favType, "error", err)
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
	if err := rows.Err(); err != nil {
		return nil, err
	}

	c.cached = &ids
	return c.cached, nil
}

func (c *LocalFavouriteCache) Bust() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cached = nil
}
