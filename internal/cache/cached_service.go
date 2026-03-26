package cache

import (
	"context"
	"log/slog"
	"strings"
	"sync"

	"github.com/RayZ3R0/sonami-gtk/internal/localdb"
	"github.com/RayZ3R0/sonami-gtk/internal/services/tidal"
	"github.com/RayZ3R0/sonami-gtk/internal/services/tidal/openapi"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
)

// CachedService wraps a tidal.Tidal service with caching
type CachedService struct {
	tidalService    *tidal.Tidal
	cache           *MetadataManager
	imagePrefetcher *ImagePrefetcher
	logger          *slog.Logger
}

// NewCachedService creates a cached wrapper around a tidal service
func NewCachedService(tidalService *tidal.Tidal) *CachedService {
	return &CachedService{
		tidalService: tidalService,
		cache:        NewMetadataManager(),
		logger:       slog.With("module", "cached-service"),
	}
}

// SetImagePrefetcher injects the image prefetcher (called after DI setup)
func (s *CachedService) SetImagePrefetcher(prefetcher *ImagePrefetcher) {
	s.imagePrefetcher = prefetcher
}

// GetAlbum retrieves album from cache or fetches from TIDAL
func (s *CachedService) GetAlbum(id string) (sonami.Album, error) {
	// Try cache first
	if album, ok := s.cache.GetAlbum(id); ok {
		return album, nil
	}

	// Fetch from TIDAL API
	s.logger.Debug("fetching from TIDAL", "type", "album", "id", id)
	albumResp, err := s.tidalService.API.OpenAPI.V2.Albums.Album(context.Background(), id, "artists.profileArt", "coverArt")
	if err != nil {
		// Check if 404/410 (item deleted)
		if isNotFoundError(err) {
			s.cache.SetAlbumError(id, err.Error())
		}
		return nil, err
	}

	// Wrap in sonami interface
	album := openapi.NewAlbum(*albumResp)

	// Store in cache
	s.cache.SetAlbum(id, album, *albumResp)

	return album, nil
}

// GetAlbumInfo retrieves album info (delegates to GetAlbum for caching)
func (s *CachedService) GetAlbumInfo(id string) (sonami.AlbumInfo, error) {
	album, err := s.GetAlbum(id)
	if err != nil {
		return nil, err
	}
	return album, nil
}

// GetArtist retrieves artist from cache or fetches from TIDAL
func (s *CachedService) GetArtist(id string) (sonami.Artist, error) {
	// Try cache first
	if artist, ok := s.cache.GetArtist(id); ok {
		return artist, nil
	}

	// Fetch from TIDAL API
	s.logger.Debug("fetching from TIDAL", "type", "artist", "id", id)
	artistResp, err := s.tidalService.API.V2.Artist.Artist(context.Background(), id)
	if err != nil {
		if isNotFoundError(err) {
			s.cache.SetArtistError(id, err.Error())
		}
		return nil, err
	}

	// Wrap using the underlying tidal service method
	artist, err := s.tidalService.GetArtist(id)
	if err != nil {
		return nil, err
	}

	s.cache.SetArtist(id, artist, *artistResp)

	return artist, nil
}

// GetArtistInfo retrieves artist info (delegates for caching)
func (s *CachedService) GetArtistInfo(id string) (sonami.ArtistInfo, error) {
	artist, err := s.GetArtist(id)
	if err != nil {
		return nil, err
	}
	return artist, nil
}

// GetPlaylist retrieves playlist from cache or fetches from TIDAL
func (s *CachedService) GetPlaylist(id string) (sonami.Playlist, error) {
	// Try cache first
	if playlist, ok := s.cache.GetPlaylist(id); ok {
		return playlist, nil
	}

	// Fetch from TIDAL API
	s.logger.Debug("fetching from TIDAL", "type", "playlist", "id", id)
	playlistResp, err := s.tidalService.API.OpenAPI.V2.Playlists.Playlist(context.Background(), id, "coverArt", "ownerProfiles.profileArt")
	if err != nil {
		if isNotFoundError(err) {
			s.cache.SetPlaylistError(id, err.Error())
		}
		return nil, err
	}

	// Wrap and cache
	playlist := openapi.NewPlaylist(*playlistResp)
	s.cache.SetPlaylist(id, playlist, *playlistResp)

	return playlist, nil
}

// GetPlaylistInfo retrieves playlist info (delegates for caching)
func (s *CachedService) GetPlaylistInfo(id string) (sonami.PlaylistInfo, error) {
	playlist, err := s.GetPlaylist(id)
	if err != nil {
		return nil, err
	}
	return playlist, nil
}

// GetTrack retrieves track from cache or fetches from TIDAL
func (s *CachedService) GetTrack(id string) (sonami.Track, error) {
	// Try cache first
	if track, ok := s.cache.GetTrack(id); ok {
		return track, nil
	}

	// Fetch from TIDAL API
	s.logger.Debug("fetching from TIDAL", "type", "track", "id", id)
	trackResp, err := s.tidalService.API.OpenAPI.V2.Tracks.Track(context.Background(), id, "albums.coverArt", "artists.profileArt")
	if err != nil {
		if isNotFoundError(err) {
			s.cache.SetTrackError(id, err.Error())
		}
		return nil, err
	}

	// Wrap and cache
	track := openapi.NewTrack(*trackResp)
	s.cache.SetTrack(id, track, *trackResp)

	return track, nil
}

// GetAlbumTracks delegates to underlying service (no caching for paginators)
func (s *CachedService) GetAlbumTracks(id string) (sonami.Paginator[sonami.Track], error) {
	return s.tidalService.GetAlbumTracks(id)
}

// GetPlaylistTracks delegates to underlying service (no caching for paginators)
func (s *CachedService) GetPlaylistTracks(id string) (sonami.Paginator[sonami.Track], error) {
	return s.tidalService.GetPlaylistTracks(id)
}

// BatchGet methods for efficient bulk fetching

// GetAlbumBatch fetches multiple albums efficiently
func (s *CachedService) GetAlbumBatch(ids []string) []sonami.Album {
	albums := make([]sonami.Album, 0, len(ids))
	uncachedIDs := make([]string, 0)

	// Check cache for each
	for _, id := range ids {
		if album, ok := s.cache.GetAlbum(id); ok {
			albums = append(albums, album)
		} else {
			uncachedIDs = append(uncachedIDs, id)
		}
	}

	// Fetch uncached concurrently (max 8 in parallel)
	if len(uncachedIDs) > 0 {
		fetched := s.fetchAlbumsConcurrent(uncachedIDs)
		albums = append(albums, fetched...)
	}

	// Prefetch cover images in background
	if s.imagePrefetcher != nil {
		coverURLs := make([]string, 0, len(albums))
		for _, album := range albums {
			if url := album.Cover(172); url != "" {
				coverURLs = append(coverURLs, url)
			}
		}
		s.imagePrefetcher.Prefetch(coverURLs)
	}

	return albums
}

func (s *CachedService) fetchAlbumsConcurrent(ids []string) []sonami.Album {
	type result struct {
		idx   int
		album sonami.Album
	}

	results := make([]result, 0, len(ids))
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 8) // Max 8 concurrent requests

	for i, id := range ids {
		wg.Add(1)
		sem <- struct{}{}
		go func(idx int, id string) {
			defer wg.Done()
			defer func() { <-sem }()

			album, err := s.GetAlbum(id)
			if err != nil {
				s.logger.Debug("failed to fetch album", "id", id, "error", err)
				return
			}

			mu.Lock()
			results = append(results, result{idx, album})
			mu.Unlock()
		}(i, id)
	}
	wg.Wait()

	// Reconstruct in order
	albums := make([]sonami.Album, len(ids))
	filled := make([]bool, len(ids))
	for _, r := range results {
		albums[r.idx] = r.album
		filled[r.idx] = true
	}

	ordered := make([]sonami.Album, 0, len(ids))
	for i, a := range albums {
		if filled[i] {
			ordered = append(ordered, a)
		}
	}
	return ordered
}

// GetArtistBatch fetches multiple artists efficiently
func (s *CachedService) GetArtistBatch(ids []string) []sonami.Artist {
	artists := make([]sonami.Artist, 0, len(ids))
	uncachedIDs := make([]string, 0)

	for _, id := range ids {
		if artist, ok := s.cache.GetArtist(id); ok {
			artists = append(artists, artist)
		} else {
			uncachedIDs = append(uncachedIDs, id)
		}
	}

	if len(uncachedIDs) > 0 {
		fetched := s.fetchArtistsConcurrent(uncachedIDs)
		artists = append(artists, fetched...)
	}

	// Prefetch cover images in background
	if s.imagePrefetcher != nil {
		coverURLs := make([]string, 0, len(artists))
		for _, artist := range artists {
			if url := artist.Cover(172); url != "" {
				coverURLs = append(coverURLs, url)
			}
		}
		s.imagePrefetcher.Prefetch(coverURLs)
	}

	return artists
}

func (s *CachedService) fetchArtistsConcurrent(ids []string) []sonami.Artist {
	type result struct {
		idx    int
		artist sonami.Artist
	}

	results := make([]result, 0, len(ids))
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 8)

	for i, id := range ids {
		wg.Add(1)
		sem <- struct{}{}
		go func(idx int, id string) {
			defer wg.Done()
			defer func() { <-sem }()

			artist, err := s.GetArtist(id)
			if err != nil {
				s.logger.Debug("failed to fetch artist", "id", id, "error", err)
				return
			}

			mu.Lock()
			results = append(results, result{idx, artist})
			mu.Unlock()
		}(i, id)
	}
	wg.Wait()

	artists := make([]sonami.Artist, len(ids))
	filled := make([]bool, len(ids))
	for _, r := range results {
		artists[r.idx] = r.artist
		filled[r.idx] = true
	}

	ordered := make([]sonami.Artist, 0, len(ids))
	for i, a := range artists {
		if filled[i] {
			ordered = append(ordered, a)
		}
	}
	return ordered
}

// GetPlaylistBatch fetches multiple playlists efficiently
func (s *CachedService) GetPlaylistBatch(ids []string) []sonami.Playlist {
	playlists := make([]sonami.Playlist, 0, len(ids))
	uncachedIDs := make([]string, 0)

	for _, id := range ids {
		if playlist, ok := s.cache.GetPlaylist(id); ok {
			playlists = append(playlists, playlist)
		} else {
			uncachedIDs = append(uncachedIDs, id)
		}
	}

	if len(uncachedIDs) > 0 {
		fetched := s.fetchPlaylistsConcurrent(uncachedIDs)
		playlists = append(playlists, fetched...)
	}

	// Prefetch cover images in background
	if s.imagePrefetcher != nil {
		coverURLs := make([]string, 0, len(playlists))
		for _, playlist := range playlists {
			if url := playlist.Cover(172); url != "" {
				coverURLs = append(coverURLs, url)
			}
		}
		s.imagePrefetcher.Prefetch(coverURLs)
	}

	return playlists
}

func (s *CachedService) fetchPlaylistsConcurrent(ids []string) []sonami.Playlist {
	type result struct {
		idx      int
		playlist sonami.Playlist
	}

	results := make([]result, 0, len(ids))
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 8)

	for i, id := range ids {
		wg.Add(1)
		sem <- struct{}{}
		go func(idx int, id string) {
			defer wg.Done()
			defer func() { <-sem }()

			playlist, err := s.GetPlaylist(id)
			if err != nil {
				s.logger.Debug("failed to fetch playlist", "id", id, "error", err)
				return
			}

			mu.Lock()
			results = append(results, result{idx, playlist})
			mu.Unlock()
		}(i, id)
	}
	wg.Wait()

	playlists := make([]sonami.Playlist, len(ids))
	filled := make([]bool, len(ids))
	for _, r := range results {
		playlists[r.idx] = r.playlist
		filled[r.idx] = true
	}

	ordered := make([]sonami.Playlist, 0, len(ids))
	for i, p := range playlists {
		if filled[i] {
			ordered = append(ordered, p)
		}
	}
	return ordered
}

// GetTrackBatch fetches multiple tracks efficiently
func (s *CachedService) GetTrackBatch(ids []string) []sonami.Track {
	tracks := make([]sonami.Track, 0, len(ids))
	uncachedIDs := make([]string, 0)

	for _, id := range ids {
		if track, ok := s.cache.GetTrack(id); ok {
			tracks = append(tracks, track)
		} else {
			uncachedIDs = append(uncachedIDs, id)
		}
	}

	if len(uncachedIDs) > 0 {
		fetched := s.fetchTracksConcurrent(uncachedIDs)
		tracks = append(tracks, fetched...)
	}

	// Prefetch cover images in background
	if s.imagePrefetcher != nil {
		coverURLs := make([]string, 0, len(tracks))
		for _, track := range tracks {
			if url := track.Cover(172); url != "" {
				coverURLs = append(coverURLs, url)
			}
		}
		s.imagePrefetcher.Prefetch(coverURLs)
	}

	return tracks
}

func (s *CachedService) fetchTracksConcurrent(ids []string) []sonami.Track {
	type result struct {
		idx   int
		track sonami.Track
	}

	results := make([]result, 0, len(ids))
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 8)

	for i, id := range ids {
		wg.Add(1)
		sem <- struct{}{}
		go func(idx int, id string) {
			defer wg.Done()
			defer func() { <-sem }()

			track, err := s.GetTrack(id)
			if err != nil {
				s.logger.Debug("failed to fetch track", "id", id, "error", err)
				return
			}

			mu.Lock()
			results = append(results, result{idx, track})
			mu.Unlock()
		}(i, id)
	}
	wg.Wait()

	tracks := make([]sonami.Track, len(ids))
	filled := make([]bool, len(ids))
	for _, r := range results {
		tracks[r.idx] = r.track
		filled[r.idx] = true
	}

	ordered := make([]sonami.Track, 0, len(ids))
	for i, t := range tracks {
		if filled[i] {
			ordered = append(ordered, t)
		}
	}
	return ordered
}

// GetCache returns the underlying cache manager (for sync operations)
func (s *CachedService) GetCache() *MetadataManager {
	return s.cache
}

// ClearMetadataCache clears all metadata from L1 and L2 caches
func (s *CachedService) ClearMetadataCache() error {
	return s.cache.ClearAll()
}

// isNotFoundError checks if error is HTTP 404 or 410
func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "404") ||
		strings.Contains(errStr, "410") ||
		strings.Contains(errStr, "Not Found") ||
		strings.Contains(errStr, "Gone")
}

// RegisterFavouriteHook registers this cached service to prefetch newly favourited items.
// This should be called once after the CachedService is created.
func (s *CachedService) RegisterFavouriteHook() {
	localdb.RegisterFavouriteAddHook(func(favType localdb.FavouriteType, id string) {
		s.logger.Debug("prefetching newly favourited item", "type", favType, "id", id)

		switch favType {
		case localdb.FavouriteAlbum:
			s.GetAlbum(id)
		case localdb.FavouriteArtist:
			s.GetArtist(id)
		case localdb.FavouritePlaylist:
			s.GetPlaylist(id)
		case localdb.FavouriteTrack:
			s.GetTrack(id)
			// FavouriteMix - no caching for mixes currently
		}
	})
}
