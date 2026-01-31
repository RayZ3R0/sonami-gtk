package player

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tonearm/pkg/tidalapi/models/v2"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/pagination"
	"github.com/go-gst/go-gst/gst"
	"github.com/infinytum/injector"
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
}

func setLoadingState() {
	PlaybackStateChanged.Notify(func(oldValue *PlaybackState) *PlaybackState {
		oldValue.Loading = true
		return oldValue
	})
}

func unsetLoadingState() {
	PlaybackStateChanged.Notify(func(oldValue *PlaybackState) *PlaybackState {
		oldValue.Loading = false
		return oldValue
	})
}

func AddTrackToUserQueue(trackId string) error {
	track, err := resolveTrack(trackId)
	if err != nil {
		return err
	}
	UserQueue.Append(track)

	// If we added a song to the queue and nothing is playing, the user likely wants to start playing the queue
	if PlaybackStateChanged.CurrentValue().Status == PlaybackStatusStopped {
		logger.Info("no track is currently playing, immediately playing track", "track_id", trackId)
		Next()
		return nil
	}
	return nil
}

func PlayTrack(trackId string) error {
	setLoadingState()
	track, err := resolveTrack(trackId)
	if err != nil {
		unsetLoadingState()
		return err
	}

	clearQueues()
	err = playTrack(track)
	if err != nil {
		return err
	}

	SourceChanged.Notify(func(oldValue *Source) *Source {
		return &Source{
			CoverURL: TrackChanged.CurrentValue().CoverURL,
			Title:    TrackChanged.CurrentValue().Title,
			Route:    fmt.Sprintf("album/%s", TrackChanged.CurrentValue().ID),
		}
	})

	history.Push(&HistoryEntry{
		TrackID: track.Data.ID,
	})

	return nil
}

func PlayAlbum(albumId string, shuffle bool, position int) error {
	setLoadingState()
	tidal, err := injector.Inject[*tidalapi.TidalAPI]()
	if err != nil {
		unsetLoadingState()
		return err
	}

	paginator := pagination.NewPaginator(tidal.OpenAPI.V2.Albums, albumId, func(items *openapi.Response[[]openapi.Relationship]) []openapi.Track {
		return items.Included.Tracks(items.Data...)
	}, "items", "items.artists", "items.albums.coverArt")

	tracks, err := paginator.GetAll()
	if err != nil {
		unsetLoadingState()
		return err
	}

	if err := PlayTracklist(tracks, shuffle, position); err != nil {
		unsetLoadingState()
		return err
	}

	SourceChanged.Notify(func(oldValue *Source) *Source {
		return &Source{
			CoverURL: TrackChanged.CurrentValue().CoverURL,
			Title:    TrackChanged.CurrentValue().Albums[0].Data.Attributes.Title,
			Route:    fmt.Sprintf("album/%s", albumId),
		}
	})

	return nil
}

func PlayArtistTopSongs(artistId string, shuffle bool, position int) error {
	setLoadingState()
	tidal, err := injector.Inject[*tidalapi.TidalAPI]()
	if err != nil {
		unsetLoadingState()
		return err
	}

	artist, err := tidal.V2.Artist.Artist(context.Background(), artistId)
	if err != nil {
		unsetLoadingState()
		return err
	}

	var module v2.PageItem
	for _, item := range artist.Items {
		if item.ModuleID == "ARTIST_TOP_TRACKS" {
			module = item
			break
		}
	}

	var topTracks []openapi.Track

	for _, LegacyTopTrackItem := range module.Items {
		if LegacyTopTrackItem.Type == v2.ItemTypeTrack {
			topTrack, _ := resolveTrack(strconv.Itoa(LegacyTopTrackItem.Data.Track.ID))
			topTracks = append(topTracks, *topTrack)
		}
	}
	PlayTracklist(topTracks, shuffle, position)

	return nil
}

func PlayPlaylist(playlistId string, shuffle bool, position int) error {
	setLoadingState()
	tidal, err := injector.Inject[*tidalapi.TidalAPI]()
	if err != nil {
		unsetLoadingState()
		return err
	}

	paginator := pagination.NewPaginator(tidal.OpenAPI.V2.Playlists, playlistId, func(items *openapi.Response[[]openapi.Relationship]) []openapi.Track {
		return items.Included.Tracks(items.Data...)
	}, "items", "items.artists", "items.albums.coverArt")

	tracks, err := paginator.GetAll()
	if err != nil {
		unsetLoadingState()
		return err
	}

	if err := PlayTracklist(tracks, shuffle, position); err != nil {
		unsetLoadingState()
		return err
	}

	playlist, err := tidal.OpenAPI.V2.Playlists.Playlist(context.Background(), playlistId, "coverArt")
	if err != nil {
		unsetLoadingState()
		return err
	}

	SourceChanged.Notify(func(oldValue *Source) *Source {
		return &Source{
			CoverURL: playlist.Included.PlainArtworks(playlist.Data.Relationships.CoverArt.Data...).AtLeast(80),
			Title:    playlist.Data.Attributes.Name,
			Route:    fmt.Sprintf("playlist/%s", playlistId),
		}
	})

	return nil
}

func PlayTrackRadio(trackId string, skipSelf bool) error {
	setLoadingState()
	tidal, err := injector.Inject[*tidalapi.TidalAPI]()
	if err != nil {
		unsetLoadingState()
		return err
	}

	trackIdInt, err := strconv.Atoi(trackId)
	if err != nil {
		logger.Error("failed to parse track id", "error", err)
		unsetLoadingState()
		return err
	}

	mix, err := tidal.V1.Tracks.Mix(context.Background(), trackIdInt)
	if err != nil {
		logger.Error("failed to retrieve mix", "error", err)
		unsetLoadingState()
		return err
	}

	position := 0
	if skipSelf {
		position = 1
	}

	return PlayPlaylist(mix.ID, false, position)
}

func PlayTracklist(tracks []openapi.Track, shuffle bool, startAt int) error {
	clearQueues()
	TrackChanged.Notify(func(oldValue *Track) *Track {
		return nil
	})

	trackPointers := make([]*openapi.Track, len(tracks))
	for i, track := range tracks {
		trackPointers[i] = &track
	}
	BaseQueue.Set(trackPointers)

	if startAt > 0 {
		BaseQueue.Skip(startAt)
	}

	if shuffle {
		SetShuffle(true)
	}

	nextTrack := getNextTrackFromQueue(false)
	if nextTrack != nil {
		logger.Info("playing next track", "track_id", nextTrack.Data.ID)
		playTrack(nextTrack)
		history.Push(&HistoryEntry{
			TrackID: nextTrack.Data.ID,
		})
	} else {
		unsetLoadingState()
	}
	return nil
}
