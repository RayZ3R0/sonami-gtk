package routes

import (
	"sync"

	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/internal/router"
	"github.com/RayZ3R0/sonami-gtk/internal/state"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components/horizontal_list"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components/media_card"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components/tracklist"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/routes/my_collection"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
	"github.com/infinytum/injector"
)

func init() {
	router.Register("my-collection", MyCollection)
	router.Register("my-collection/albums", my_collection.Albums)
	router.Register("my-collection/artists", my_collection.Artists)
	router.Register("my-collection/playlists", my_collection.Playlists)
	router.Register("my-collection/tracks", my_collection.Tracks)
}

// fetchN fetches up to n items from ids concurrently (max 8 in-flight).
// Results preserve the order of ids; items that fail to load are skipped.
func fetchN[T any](ids []string, n int, fetch func(string) (T, error)) []T {
	if n > 0 && len(ids) > n {
		ids = ids[:n]
	}

	type indexed struct {
		i   int
		val T
	}

	results := make([]indexed, 0, len(ids))
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 8)

	for i, id := range ids {
		i, id := i, id
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			val, err := fetch(id)
			if err != nil {
				return
			}
			mu.Lock()
			results = append(results, indexed{i, val})
			mu.Unlock()
		}()
	}
	wg.Wait()

	// Sort by original index to preserve added_at order from the DB.
	out := make([]T, len(ids))
	filled := make([]bool, len(ids))
	for _, r := range results {
		out[r.i] = r.val
		filled[r.i] = true
	}
	ordered := make([]T, 0, len(ids))
	for i, v := range out {
		if filled[i] {
			ordered = append(ordered, v)
		}
	}
	return ordered
}

func MyCollection() *router.Response {
	service, err := injector.Inject[sonami.Service]()
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

	// Fetch up to 8 items per section in parallel across sections.
	var (
		albums    []sonami.Album
		artists   []sonami.Artist
		playlists []sonami.Playlist
		tracks    []sonami.Track
		wg        sync.WaitGroup
	)

	wg.Add(4)
	go func() {
		defer wg.Done()
		albums = fetchN(*albumIDs, 8, service.GetAlbum)
	}()
	go func() {
		defer wg.Done()
		artists = fetchN(*artistIDs, 8, service.GetArtist)
	}()
	go func() {
		defer wg.Done()
		playlists = fetchN(*playlistIDs, 8, service.GetPlaylist)
	}()
	go func() {
		defer wg.Done()
		tracks = fetchN(*trackIDs, 8, service.GetTrack)
	}()
	wg.Wait()

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

	// Playlists section.
	if len(playlists) > 0 {
		playlistList := horizontal_list.NewHorizontalList(gettext.Get("Playlists")).
			SetPageMargin(40)
		if len(*playlistIDs) > 8 {
			playlistList.SetViewAllRoute("my-collection/playlists")
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
	if len(albums) == 0 && len(artists) == 0 && len(playlists) == 0 && len(tracks) == 0 {
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
