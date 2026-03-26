package routes

import (
	"context"

	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/RayZ3R0/sonami-gtk/internal/cache"
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/internal/localdb"
	"github.com/RayZ3R0/sonami-gtk/internal/router"
	"github.com/RayZ3R0/sonami-gtk/internal/state"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components/horizontal_list"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components/media_card"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components/tracklist"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/routes/my_collection"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	"github.com/infinytum/injector"
)

func init() {
	router.Register("my-collection", MyCollection)
	router.Register("my-collection/albums", my_collection.Albums)
	router.Register("my-collection/artists", my_collection.Artists)
	router.Register("my-collection/playlists", my_collection.Playlists)
	router.Register("my-collection/tracks", my_collection.Tracks)
}

func MyCollection() *router.Response {
	cachedService, err := injector.Inject[*cache.CachedService]()
	if err != nil {
		return router.FromError(gettext.Get("My Collection"), err)
	}

	albumIDs, err := state.AlbumsCache.Get()
	if err != nil {
		return router.FromError(gettext.Get("My Collection"), err)
	}
	artistIDs, err := state.ArtistsCache.Get()
	if err != nil {
		return router.FromError(gettext.Get("My Collection"), err)
	}
	playlistIDs, err := state.PlaylistsCache.Get()
	if err != nil {
		return router.FromError(gettext.Get("My Collection"), err)
	}
	trackIDs, err := state.TracksCache.Get()
	if err != nil {
		return router.FromError(gettext.Get("My Collection"), err)
	}

	// Local playlists are a fast DB-only operation — fetch synchronously before goroutines.
	localPlaylists, _ := localdb.GetAllPlaylists()
	maxTidalPlaylists := 8 - len(localPlaylists)
	if maxTidalPlaylists < 0 {
		maxTidalPlaylists = 0
	}

	// Limit IDs to first 8 per section
	limitedAlbumIDs := *albumIDs
	if len(limitedAlbumIDs) > 8 {
		limitedAlbumIDs = limitedAlbumIDs[:8]
	}
	limitedArtistIDs := *artistIDs
	if len(limitedArtistIDs) > 8 {
		limitedArtistIDs = limitedArtistIDs[:8]
	}
	limitedPlaylistIDs := *playlistIDs
	if len(limitedPlaylistIDs) > maxTidalPlaylists {
		limitedPlaylistIDs = limitedPlaylistIDs[:maxTidalPlaylists]
	}
	limitedTrackIDs := *trackIDs
	if len(limitedTrackIDs) > 8 {
		limitedTrackIDs = limitedTrackIDs[:8]
	}

	// Use batch operations for efficient cache-aware fetching
	albums := cachedService.GetAlbumBatch(limitedAlbumIDs)
	artists := cachedService.GetArtistBatch(limitedArtistIDs)
	playlists := cachedService.GetPlaylistBatch(limitedPlaylistIDs)
	tracks := cachedService.GetTrackBatch(limitedTrackIDs)

	// Trigger background sync for stale entries (non-blocking)
	if syncManager, err := injector.Inject[*cache.SyncManager](); err == nil {
		go syncManager.SyncStaleEntries(context.Background())
	}

	// Prefetch all liked tracks and local playlist tracks in background
	go func() {
		// Prefetch ALL liked tracks (not just the displayed 8)
		if len(*trackIDs) > 0 {
			cachedService.GetTrackBatch(*trackIDs)
		}

		// Prefetch tracks for all local playlists
		for _, lp := range localPlaylists {
			if trackIDs, err := localdb.GetPlaylistTrackIDs(lp.ID); err == nil && len(trackIDs) > 0 {
				cachedService.GetTrackBatch(trackIDs)
			}
		}
	}()

	body := VStack().Spacing(25).VMargin(20)

	// Artists section.
	if len(artists) > 0 {
		artistList := horizontal_list.NewHorizontalList(gettext.Get("Artists")).
			SetPageMargin(40)
		if len(*artistIDs) > 8 {
			artistList.SetViewAllRoute("my-collection/artists")
		}
		for _, artist := range artists {
			artistList.Append(media_card.NewArtist(artist))
		}
		body = body.Append(artistList)
	}

	// Albums section.
	if len(albums) > 0 {
		albumList := horizontal_list.NewHorizontalList(gettext.Get("Albums")).
			SetPageMargin(40)
		if len(*albumIDs) > 8 {
			albumList.SetViewAllRoute("my-collection/albums")
		}
		for _, album := range albums {
			albumList.Append(media_card.NewAlbum(album))
		}
		body = body.Append(albumList)
	}

	// Playlists section (local user-created playlists first, then Tidal-saved).
	hasLocalPlaylists := len(localPlaylists) > 0
	hasPlaylists := hasLocalPlaylists || len(playlists) > 0
	if hasPlaylists {
		playlistList := horizontal_list.NewHorizontalList(gettext.Get("Playlists")).
			SetPageMargin(40)
		totalPlaylists := len(localPlaylists) + len(*playlistIDs)
		if totalPlaylists > 8 {
			playlistList.SetViewAllRoute("my-collection/playlists")
		}
		for _, lp := range localPlaylists {
			playlistList.Append(media_card.NewLocalPlaylist(lp))
		}
		for _, playlist := range playlists {
			playlistList.Append(media_card.NewPlaylist(playlist))
		}
		body = body.Append(playlistList)
	}

	// Tracks section.
	if len(tracks) > 0 {
		trackList := tracklist.NewTrackList(
			tracklist.CoverColumn, tracklist.TitleAlbumColumn,
			tracklist.ArtistsColumn,
			tracklist.DurationColumn, tracklist.ControlsColumn,
		)
		for _, track := range tracks {
			trackList.AddTrack(track)
		}
		tracksSection := VStack(
			components.NewRowTitle().
				SetTitle(gettext.Get("Tracks")).
				SetViewAllRoute("my-collection/tracks"),
			trackList,
		).HMargin(40)
		body = body.Append(tracksSection)
	}

	// Empty state: all sections are empty.
	if len(albums) == 0 && len(artists) == 0 && !hasPlaylists && len(tracks) == 0 {
		return &router.Response{
			PageTitle: gettext.Get("My Collection"),
			View: StatusPage().
				IconName("heart-outline-thick-symbolic").
				Title(gettext.Get("Nothing Here Yet")).
				Description(gettext.Get("Tap the heart icon on any track, album, artist, or playlist to save it here")),
		}
	}

	return &router.Response{
		PageTitle: gettext.Get("My Collection"),
		View: ScrolledWindow().
			Child(
				components.MainContent(
					body.
						Append(Spacer()).
						VAlign(gtk.AlignStartValue),
				),
			).
			Policy(gtk.PolicyNeverValue, gtk.PolicyAutomaticValue),
	}
}
