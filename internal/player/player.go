package player

import (
	"log/slog"

	"codeberg.org/dergs/tidalwave/internal/settings"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/pagination"
	"github.com/go-gst/go-gst/gst"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gobject"
)

var (
	logger  = slog.With("module", "player")
	playbin *gst.Element
)

func init() {
	gst.Init(nil)
	pb, err := gst.NewElement("playbin")
	if err != nil {
		panic(err)
	}
	playbin = pb
	playbin.GetBus().AddWatch(onBusMessage)
	playbin.Connect("notify::volume", onVolumeChange)
	playbin.Connect("about-to-finish", onAboutToFinish)
	settings.PlayerSettings().BindVolume(gobject.ObjectNewFromInternalPtr(uintptr(playbin.BaseObject().Unsafe())), "volume")
}

func AddTrackToUserQueue(trackId string) {
	UserQueue.AddTrackID(trackId, false)

	// If we added a song to the queue and nothing is playing, the user likely wants to start playing the queue
	if PlaybackStateChanged.CurrentValue().Status == PlaybackStatusStopped {
		logger.Info("no track is currently playing, immediately playing track", "track_id", trackId)
		Next()
		return
	}
}

func PlayTrack(trackId string) error {
	track, err := resolveTrack(trackId)
	if err != nil {
		return err
	}

	return playTrack(track)
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

	clearQueues()
	for _, track := range prepareTrackList(tracks, shuffle, skipUntil) {
		BaseQueue.AddTrack(&track, false)
	}

	nextTrack := getNextTrackFromQueue(false)
	if nextTrack != nil {
		logger.Info("playing next track", "track_id", nextTrack.Data.ID)
		playTrack(nextTrack)
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
	}, "items", "items.artists", "items.albums.coverArt")

	tracks, err := paginator.GetAll()
	if err != nil {
		return err
	}

	clearQueues()
	for _, track := range prepareTrackList(tracks, shuffle, skipUntil) {
		BaseQueue.AddTrack(&track, false)
	}
	nextTrack := getNextTrackFromQueue(false)
	if nextTrack != nil {
		logger.Info("playing next track", "track_id", nextTrack.Data.ID)
		playTrack(nextTrack)
	}

	return nil
}
