package player

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/RayZ3R0/sonami-gtk/internal/settings"
	"github.com/RayZ3R0/sonami-gtk/internal/signals"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi"
	v2 "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/models/v2"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
	"github.com/go-gst/go-gst/gst"
	"github.com/infinytum/injector"
)

var (
	logger  = slog.With("module", "player")
	playbin *gst.Element
)

func init() {
	gst.Init(nil)
	pb, err := gst.NewElement("playbin3")
	if err != nil {
		panic(err)
	}
	playbin = pb
	playbin.GetBus().AddWatch(onBusMessage)
	playbin.Connect("notify::volume", onVolumeChange)
	playbin.Connect("about-to-finish", onAboutToFinish)
	playbin.Connect("deep-element-added", onDeepElementAdded)
	playbin.SetProperty("buffer-size", 20*1024*1024)                         // 20 MB
	playbin.SetProperty("buffer-duration", (30 * time.Second).Nanoseconds()) // 30 seconds

	audioFilterBin, err := buildReplayGainFilterBin()
	if err != nil {
		panic(err)
	}

	emptyBin, err := buildEmptyBin()
	if err != nil {
		panic(err)
	}

	settings.Playback().ConnectNormalizeVolumeChanged(func(newVal bool) bool {
		var err error
		if newVal {
			err = playbin.Set("audio-filter", audioFilterBin)
		} else {
			err = playbin.Set("audio-filter", emptyBin)
		}
		if err != nil {
			logger.Error("Failed to set audio filter", err)
		}
		return signals.Continue
	})
}

var stateBeforeLoading gst.State

func setLoadingState() {
	_, stateBeforeLoading = playbin.GetState(gst.StatePaused, gst.ClockTimeNone)
	PlaybackStateChanged.Notify(func(oldValue *PlaybackState) *PlaybackState {
		oldValue.Loading = true
		return oldValue
	})
}

func resetLoadingState() {
	if stateBeforeLoading != gst.VoidPending {
		playbin.SetState(stateBeforeLoading)
		stateBeforeLoading = gst.VoidPending
	}

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

func AddTracklistToUserQueue(tracklist []sonami.Track) error {
	for _, track := range tracklist {
		UserQueue.Append(track)
	}

	// If we added a song to the queue and nothing is playing, the user likely wants to start playing the queue
	if PlaybackStateChanged.CurrentValue().Status == PlaybackStatusStopped && len(tracklist) > 0 {
		logger.Info("no track is currently playing, immediately playing track", "track_id", tracklist[0].ID())
		Next()
		return nil
	}
	return nil
}

func PlayTrackID(trackId string) error {
	setLoadingState()
	track, err := resolveTrack(trackId)
	if err != nil {
		resetLoadingState()
		return err
	}

	clearQueues()
	err = playTrack(track)
	if err != nil {
		return err
	}

	SourceChanged.Notify(func(oldValue sonami.PlaybackSource) sonami.PlaybackSource {
		return track
	})

	history.Push(&HistoryEntry{
		TrackID: track.ID(),
	})

	return nil
}

func PlayTrack(track sonami.Track) error {
	setLoadingState()

	clearQueues()
	err := playTrack(track)
	if err != nil {
		return err
	}

	SourceChanged.Notify(func(oldValue sonami.PlaybackSource) sonami.PlaybackSource {
		return track
	})

	history.Push(&HistoryEntry{
		TrackID: track.ID(),
	})

	return nil
}

func PlayTracklist(source sonami.PlaybackSource, tracklist []sonami.Track, shuffle bool, position int) error {
	setLoadingState()
	if _, err := playTracklist(tracklist, shuffle, position); err != nil {
		resetLoadingState()
		return err
	}

	SourceChanged.Notify(func(oldValue sonami.PlaybackSource) sonami.PlaybackSource {
		return source
	})
	return nil
}

func PlayAlbum(albumId string, shuffle bool, position int) error {
	setLoadingState()
	if err := playAlbum(albumId, shuffle, position); err != nil {
		resetLoadingState()
		return err
	}
	return nil
}

func PlayArtistTopSongs(artistId string, shuffle bool, position int) error {
	setLoadingState()
	tidal, err := injector.Inject[*tidalapi.TidalAPI]()
	if err != nil {
		resetLoadingState()
		return err
	}

	service, err := injector.Inject[sonami.Service]()
	if err != nil {
		resetLoadingState()
		return err
	}

	artistInfo, err := service.GetArtistInfo(artistId)
	if err != nil {
		resetLoadingState()
		return err
	}

	artist, err := tidal.V2.Artist.Artist(context.Background(), artistId)
	if err != nil {
		resetLoadingState()
		return err
	}

	var module v2.PageItem
	for _, item := range artist.Items {
		if item.ModuleID == "ARTIST_TOP_TRACKS" {
			module = item
			break
		}
	}

	var topTracks []sonami.Track

	for _, legacyTopTrackItem := range module.Items {
		if legacyTopTrackItem.Type == v2.ItemTypeTrack {
			topTrack, err := resolveTrack(strconv.Itoa(legacyTopTrackItem.Data.Track.ID))
			if err != nil {
				logger.Error("error while resolving Top Track item", "track_id", legacyTopTrackItem.Data.Track.ID, "message", err.Error())
				continue
			}

			topTracks = append(topTracks, topTrack)
		}
	}
	_, err = playTracklist(topTracks, shuffle, position)

	if err == nil {
		SourceChanged.Notify(func(oldValue sonami.PlaybackSource) sonami.PlaybackSource {
			return artistInfo
		})
	}

	return err
}

func PlayPlaylist(playlistId string, shuffle bool, position int) error {
	setLoadingState()
	if err := playPlaylist(playlistId, shuffle, position); err != nil {
		resetLoadingState()
		return err
	}
	return nil
}

func PlayTrackRadio(trackId string, skipSelf bool) error {
	setLoadingState()
	tidal, err := injector.Inject[*tidalapi.TidalAPI]()
	if err != nil {
		resetLoadingState()
		return err
	}

	trackIdInt, err := strconv.Atoi(trackId)
	if err != nil {
		logger.Error("failed to parse track id", "error", err)
		resetLoadingState()
		return err
	}

	mix, err := tidal.V1.Tracks.Mix(context.Background(), trackIdInt)
	if err != nil {
		logger.Error("failed to retrieve mix", "error", err)
		resetLoadingState()
		return err
	}

	position := 0
	if skipSelf {
		position = 1
	}

	return PlayPlaylist(mix.ID, false, position)
}
