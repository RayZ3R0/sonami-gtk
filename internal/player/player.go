package player

import (
	"log/slog"
	"math/rand/v2"

	"codeberg.org/dergs/tidalwave/internal/settings"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/pagination"
	"github.com/go-gst/go-gst/gst"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gobject"
)

var (
	playbin                    *gst.Element
	currentlyConfiguredTrackID int
	logger                     *slog.Logger
)

func init() {
	gst.Init(nil)
	logger = slog.With("module", "player")
	pb, err := gst.NewElement("playbin")
	if err != nil {
		panic(err)
	}
	playbin = pb

	settings.PlayerSettings().BindVolume(gobject.ObjectNewFromInternalPtr(uintptr(playbin.BaseObject().Unsafe())), "volume")
	playbin.GetBus().AddWatch(onBusMessage)
	playbin.Connect("notify::volume", onVolumeChange)
	playbin.Connect("about-to-finish", onAboutToFinish)
	onVolumeChange()
}

func PlayTrack(trackId string) error {
	BaseQueue.Clear()
	return playTrackId(trackId)
}

func PlayAlbum(albumId string, shuffle bool, skipUntil string) error {
	tidal, err := injector.Inject[*tidalapi.TidalAPI]()
	if err != nil {
		return err
	}

	paginator := pagination.NewPaginator(tidal.OpenAPI.V2.Albums, albumId, func(items *openapi.Response[[]openapi.Relationship]) []openapi.Track {
		return items.Included.Tracks(items.Data...)
	}, "items", "items.artists", "items.albums.coverArt")
	tracks, err := paginator.GetAll()

	if err != nil {
		return err
	}

	BaseQueue.Clear()

	if shuffle {
		rand.Shuffle(len(tracks), func(i, j int) {
			tracks[i], tracks[j] = tracks[j], tracks[i]
		})
	} else if skipUntil != "" {
		for i, track := range tracks {
			if track.Data.ID == skipUntil {
				tracks = tracks[i:]
				break
			}
		}
	}

	firstTrack := tracks[0]

	if err := playTrack(&firstTrack); err != nil {
		return err
	}

	for _, track := range tracks[1:] {
		BaseQueue.AddTrack(&track, false)
	}

	return nil
}

func PlayPlaylist(playlistId string, shuffle bool, skipUntil string) error {
	tidal, err := injector.Inject[*tidalapi.TidalAPI]()
	if err != nil {
		return err
	}

	paginator := pagination.NewPaginator(tidal.OpenAPI.V2.Playlists, playlistId, func(items *openapi.Response[[]openapi.Relationship]) []openapi.Track {
		return items.Included.Tracks(items.Data...)
	})
	tracks, err := paginator.GetAll()

	if err != nil {
		return err
	}

	BaseQueue.Clear()

	if shuffle {
		rand.Shuffle(len(tracks), func(i, j int) {
			tracks[i], tracks[j] = tracks[j], tracks[i]
		})
	} else if skipUntil != "" {
		for i, track := range tracks {
			if track.Data.ID == skipUntil {
				tracks = tracks[i:]
				break
			}
		}
	}

	firstTrack := tracks[0]
	if err := playTrack(&firstTrack); err != nil {
		return err
	}

	for _, track := range tracks[1:] {
		BaseQueue.AddTrack(&track, false)
	}

	return nil
}
