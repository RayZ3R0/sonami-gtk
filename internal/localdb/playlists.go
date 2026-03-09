package localdb

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"time"
)

var playlistsLogger = slog.With("module", "localdb/playlists")

// LocalPlaylist is a user-created playlist stored in the local SQLite database.
type LocalPlaylist struct {
	ID         string
	Name       string
	CoverURL   string
	TrackCount int
	CreatedAt  time.Time
}

func generateID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%016x", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}

// CreatePlaylist creates a new empty local playlist with the given name.
func CreatePlaylist(name string) (LocalPlaylist, error) {
	id := generateID()
	_, err := DB().Exec(`INSERT INTO local_playlists (id, name) VALUES (?, ?)`, id, name)
	if err != nil {
		playlistsLogger.Error("failed to create playlist", "name", name, "error", err)
		return LocalPlaylist{}, err
	}
	return LocalPlaylist{ID: id, Name: name, CreatedAt: time.Now()}, nil
}

// DeletePlaylist removes a playlist and all its associated tracks.
func DeletePlaylist(id string) error {
	if _, err := DB().Exec(`DELETE FROM local_playlist_tracks WHERE playlist_id = ?`, id); err != nil {
		playlistsLogger.Error("failed to delete playlist tracks", "id", id, "error", err)
		return err
	}
	if _, err := DB().Exec(`DELETE FROM local_playlists WHERE id = ?`, id); err != nil {
		playlistsLogger.Error("failed to delete playlist", "id", id, "error", err)
		return err
	}
	return nil
}

// RenamePlaylist updates the display name of a local playlist.
func RenamePlaylist(id, newName string) error {
	_, err := DB().Exec(`UPDATE local_playlists SET name = ? WHERE id = ?`, newName, id)
	if err != nil {
		playlistsLogger.Error("failed to rename playlist", "id", id, "error", err)
	}
	return err
}

// GetAllPlaylists returns all local playlists ordered by creation date (newest first).
func GetAllPlaylists() ([]LocalPlaylist, error) {
	rows, err := DB().Query(`
		SELECT lp.id, lp.name, lp.cover_url, COUNT(lpt.track_id) AS track_count, lp.created_at
		FROM local_playlists lp
		LEFT JOIN local_playlist_tracks lpt ON lp.id = lpt.playlist_id
		GROUP BY lp.id
		ORDER BY lp.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var playlists []LocalPlaylist
	for rows.Next() {
		var p LocalPlaylist
		var createdAtStr string
		if err := rows.Scan(&p.ID, &p.Name, &p.CoverURL, &p.TrackCount, &createdAtStr); err != nil {
			playlistsLogger.Warn("failed to scan playlist row", "error", err)
			continue
		}
		p.CreatedAt, _ = time.Parse("2006-01-02T15:04:05.000Z", createdAtStr)
		playlists = append(playlists, p)
	}
	return playlists, nil
}

// GetPlaylist returns a single local playlist by ID, or an error if not found.
func GetPlaylist(id string) (*LocalPlaylist, error) {
	var p LocalPlaylist
	var createdAtStr string
	err := DB().QueryRow(`
		SELECT lp.id, lp.name, lp.cover_url, COUNT(lpt.track_id), lp.created_at
		FROM local_playlists lp
		LEFT JOIN local_playlist_tracks lpt ON lp.id = lpt.playlist_id
		WHERE lp.id = ?
		GROUP BY lp.id
	`, id).Scan(&p.ID, &p.Name, &p.CoverURL, &p.TrackCount, &createdAtStr)
	if err != nil {
		return nil, err
	}
	p.CreatedAt, _ = time.Parse("2006-01-02T15:04:05.000Z", createdAtStr)
	return &p, nil
}

// AddTrackToPlaylist appends a track at the end of a local playlist
// and updates the playlist cover to the given coverURL (last-added wins).
// Silently ignores duplicates (INSERT OR IGNORE).
func AddTrackToPlaylist(playlistID, trackID, coverURL string) error {
	var maxPos int
	_ = DB().QueryRow(
		`SELECT COALESCE(MAX(position), -1) FROM local_playlist_tracks WHERE playlist_id = ?`,
		playlistID,
	).Scan(&maxPos)
	_, err := DB().Exec(
		`INSERT OR IGNORE INTO local_playlist_tracks (playlist_id, track_id, position) VALUES (?, ?, ?)`,
		playlistID, trackID, maxPos+1,
	)
	if err != nil {
		playlistsLogger.Error("failed to add track to playlist", "playlist_id", playlistID, "track_id", trackID, "error", err)
		return err
	}
	if coverURL != "" {
		_, _ = DB().Exec(`UPDATE local_playlists SET cover_url = ? WHERE id = ?`, coverURL, playlistID)
	}
	return nil
}

// RemoveTrackFromPlaylist removes a specific track from a local playlist.
func RemoveTrackFromPlaylist(playlistID, trackID string) error {
	_, err := DB().Exec(
		`DELETE FROM local_playlist_tracks WHERE playlist_id = ? AND track_id = ?`,
		playlistID, trackID,
	)
	if err != nil {
		playlistsLogger.Error("failed to remove track from playlist", "playlist_id", playlistID, "track_id", trackID, "error", err)
	}
	return err
}

// GetPlaylistTrackIDs returns the track IDs in a playlist, ordered by position then add time.
func GetPlaylistTrackIDs(playlistID string) ([]string, error) {
	rows, err := DB().Query(
		`SELECT track_id FROM local_playlist_tracks WHERE playlist_id = ? ORDER BY position ASC, added_at ASC`,
		playlistID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err == nil {
			ids = append(ids, id)
		}
	}
	return ids, nil
}
