package localdb

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	_ "modernc.org/sqlite"
)

const appID = "io.github.rayz3r0.SonamiGtk"

var (
	db     *sql.DB
	dbOnce sync.Once
)

// DB returns the shared database connection.
// It is lazily initialized on first call and safe for concurrent use.
func DB() *sql.DB {
	dbOnce.Do(func() {
		dbPath := filepath.Join(dataDir(), "library.db")

		var err error
		db, err = sql.Open("sqlite", dbPath+"?_journal_mode=WAL&_busy_timeout=5000")
		if err != nil {
			slog.Error("failed to open local database", "path", dbPath, "error", err)
			panic("localdb: " + err.Error())
		}

		db.SetMaxOpenConns(1)

		if err := migrate(db); err != nil {
			slog.Error("failed to run database migrations", "error", err)
			panic("localdb: migrate: " + err.Error())
		}

		slog.Info("local database ready", "path", dbPath)
	})
	return db
}

// dataDir returns the XDG data directory for the application,
// creating it if it doesn't exist.
func dataDir() string {
	base := os.Getenv("XDG_DATA_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			slog.Error("failed to get user home directory", "error", err)
			panic("localdb: " + err.Error())
		}
		base = filepath.Join(home, ".local", "share")
	}

	dir := filepath.Join(base, appID)
	if err := os.MkdirAll(dir, 0700); err != nil {
		slog.Error("failed to create data directory", "path", dir, "error", err)
		panic("localdb: " + err.Error())
	}
	return dir
}

// migrate applies the database schema.
// Uses a user_version pragma to track schema version.
func migrate(db *sql.DB) error {
	var version int
	if err := db.QueryRow("PRAGMA user_version").Scan(&version); err != nil {
		return err
	}

	migrations := []string{
		// v1: favourites table
		`CREATE TABLE IF NOT EXISTS favourites (
			type     TEXT NOT NULL,
			id       TEXT NOT NULL,
			added_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
			PRIMARY KEY (type, id)
		)`,
		// v2: local playlists
		`CREATE TABLE IF NOT EXISTS local_playlists (
			id         TEXT PRIMARY KEY,
			name       TEXT NOT NULL,
			created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
		)`,
		// v3: local playlist tracks
		`CREATE TABLE IF NOT EXISTS local_playlist_tracks (
			playlist_id TEXT NOT NULL,
			track_id    TEXT NOT NULL,
			position    INTEGER NOT NULL DEFAULT 0,
			added_at    TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
			PRIMARY KEY (playlist_id, track_id)
		)`,
		// v4: add cover_url to local playlists
		`ALTER TABLE local_playlists ADD COLUMN cover_url TEXT NOT NULL DEFAULT ''`,
		// v5: metadata cache for albums
		`CREATE TABLE IF NOT EXISTS cached_albums (
			id TEXT PRIMARY KEY,
			data BLOB NOT NULL,
			cover_url TEXT NOT NULL DEFAULT '',
			cached_at INTEGER NOT NULL,
			error TEXT DEFAULT NULL
		)`,
		// v6: metadata cache for artists
		`CREATE TABLE IF NOT EXISTS cached_artists (
			id TEXT PRIMARY KEY,
			data BLOB NOT NULL,
			cover_url TEXT NOT NULL DEFAULT '',
			cached_at INTEGER NOT NULL,
			error TEXT DEFAULT NULL
		)`,
		// v7: metadata cache for playlists
		`CREATE TABLE IF NOT EXISTS cached_playlists (
			id TEXT PRIMARY KEY,
			data BLOB NOT NULL,
			cover_url TEXT NOT NULL DEFAULT '',
			cached_at INTEGER NOT NULL,
			error TEXT DEFAULT NULL
		)`,
		// v8: metadata cache for tracks
		`CREATE TABLE IF NOT EXISTS cached_tracks (
			id TEXT PRIMARY KEY,
			data BLOB NOT NULL,
			cover_url TEXT NOT NULL DEFAULT '',
			cached_at INTEGER NOT NULL,
			error TEXT DEFAULT NULL
		)`,
		// v9: indexes for stale entry queries
		`CREATE INDEX IF NOT EXISTS idx_albums_cached_at ON cached_albums(cached_at)`,
		// v10: indexes for artists
		`CREATE INDEX IF NOT EXISTS idx_artists_cached_at ON cached_artists(cached_at)`,
		// v11: indexes for playlists
		`CREATE INDEX IF NOT EXISTS idx_playlists_cached_at ON cached_playlists(cached_at)`,
		// v12: indexes for tracks
		`CREATE INDEX IF NOT EXISTS idx_tracks_cached_at ON cached_tracks(cached_at)`,
		// v13: page feed cache (for home, explore, etc.)
		`CREATE TABLE IF NOT EXISTS cached_page_feeds (
			id TEXT PRIMARY KEY,
			data BLOB NOT NULL,
			cached_at INTEGER NOT NULL
		)`,
	}

	for i := version; i < len(migrations); i++ {
		if _, err := db.Exec(migrations[i]); err != nil {
			return err
		}
	}

	if version < len(migrations) {
		_, err := db.Exec(fmt.Sprintf("PRAGMA user_version = %d", len(migrations)))
		return err
	}
	return nil
}
