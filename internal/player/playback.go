package player

import (
	"context"
	"errors"
	"slices"
	"strconv"

	"codeberg.org/dergs/tidalwave/internal/notifications"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	v1 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v1"
	tracksv1 "codeberg.org/dergs/tidalwave/pkg/tidalapi/v1/tracks"
	"github.com/go-gst/go-gst/gst"
	"github.com/infinytum/injector"
)

func playTrack(track *openapi.Track) error {
	tidal, err := injector.Inject[*tidalapi.TidalAPI]()
	if err != nil {
		return err
	}

	TrackChanged.Notify(func(oldState *Track) *Track {
		trackInfo := &Track{
			Albums:   []openapi.Album{},
			Artists:  []openapi.ArtistData{},
			CoverURL: "",
			Duration: track.Data.Attributes.Duration.Duration,
			ID:       track.Data.ID,
			Title:    track.Data.Attributes.Title,
		}

		for _, album := range track.Included.Albums(track.Data.Relationships.Albums.Data...) {
			trackInfo.Albums = append(trackInfo.Albums, album)
			for _, artwork := range album.Included.PlainArtworks(album.Data.Relationships.CoverArt.Data...) {
				trackInfo.CoverURL = artwork.Attributes.Files.AtLeast(320).Href
			}
		}

		for _, artist := range track.Included.PlainArtists(track.Data.Relationships.Artists.Data...) {
			trackInfo.Artists = append(trackInfo.Artists, artist)
		}

		return trackInfo
	})

	PlaybackStateChanged.Notify(func(oldValue *PlaybackState) *PlaybackState {
		newState := *oldValue
		newState.Duration = track.Data.Attributes.Duration.Duration
		return &newState
	})

	history.Push(&HistoryEntry{
		TrackID: track.Data.ID,
	})

	if !slices.Contains(track.Data.Attributes.Availability, openapi.TrackAvailabilityStream) {
		notifications.OnToast.Notify("Track not available for streaming, skipping to next track")
		Next()
		return errors.New("track not available for streaming")
	}

	if strconv.Itoa(currentlyEnqueuedTrackID) != track.Data.ID {
		logger.Debug("fetching playback info for track", "track_id", track.Data.ID)
		playbackInfo, err := tidal.V1.Tracks.PlaybackInfo(context.Background(), track.Data.ID, tracksv1.PlaybackInfoOptions{})
		if err != nil {
			logger.Error("unable to fetch playback info for track", "error", err)
			return err
		}
		return play(playbackInfo)
	}
	logger.Debug("gapless playback detected, not enqueueing track again")
	unsetLoadingState()
	return nil
}

func play(playbackInfo *v1.PlaybackInfo) error {
	// Inform the UI about the track quality we got from TIDAL.
	PlaybackQualityChanged.Notify(func(oldValue v1.AudioQuality) v1.AudioQuality {
		return playbackInfo.AudioQuality
	})

	// Free up resources taken up by previous stream
	playbin.SetState(gst.StateNull)
	playbin.SetArg("uri", "")

	PlaybackStateChanged.Notify(func(oldValue *PlaybackState) *PlaybackState {
		newState := *oldValue
		newState.Status = PlaybackStatusBuffering
		return &newState
	})

	if err := enqueue(playbackInfo); err != nil {
		return err
	}

	playbin.SetState(gst.StatePlaying)
	return nil
}
