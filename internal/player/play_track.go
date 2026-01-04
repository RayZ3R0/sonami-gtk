package player

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"codeberg.org/dergs/tidalwave/internal/notifications"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	v1 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v1"
	tracksv1 "codeberg.org/dergs/tidalwave/pkg/tidalapi/v1/tracks"
	"github.com/go-gst/go-gst/gst"
	"github.com/infinytum/injector"
)

func playTrackId(trackId string) error {
	tidal, err := injector.Inject[*tidalapi.TidalAPI]()
	if err != nil {
		return err
	}

	track, err := tidal.OpenAPI.V2.Tracks.Track(context.Background(), trackId, "albums.coverArt", "artists")
	if err != nil {
		return err
	}

	return playTrack(track)
}

func playTrack(track *openapi.Track) error {
	tidal, err := injector.Inject[*tidalapi.TidalAPI]()

	OnTrackChanged.Notify(func(trackInfo *TrackInformation) {
		trackInfo.Albums = []openapi.Album{}
		trackInfo.Artists = []openapi.ArtistData{}
		trackInfo.CoverURL = ""
		trackInfo.Duration = track.Data.Attributes.Duration.Duration
		trackInfo.ID = track.Data.ID
		trackInfo.Title = track.Data.Attributes.Title

		for _, artist := range track.Included.PlainArtists(track.Data.Relationships.Artists.Data...) {
			trackInfo.Artists = append(trackInfo.Artists, artist)
		}

		for _, album := range track.Included.Albums(track.Data.Relationships.Albums.Data...) {
			trackInfo.Albums = append(trackInfo.Albums, album)
			for _, artwork := range album.Included.PlainArtworks(album.Data.Relationships.CoverArt.Data...) {
				trackInfo.CoverURL = artwork.Attributes.Files.AtLeast(320).Href
			}
		}
	})

	if !slices.Contains(track.Data.Attributes.Availability, openapi.TrackAvailabilityStream) {
		notifications.OnToast.Notify("Track not available for streaming, skipping to next track")
		Next()
		return errors.New("track not available for streaming")
	}

	playbackInfo, err := tidal.V1.Tracks.PlaybackInfo(context.Background(), track.Data.ID, tracksv1.PlaybackInfoOptions{})
	if err != nil {
		logger.Error("unable to fetch playback info for track", "error", err)
		return err
	}

	return play(playbackInfo)
}

func play(playbackInfo *v1.PlaybackInfo) error {
	// Inform the UI about the track quality we got from TIDAL.
	OnPlaybackQualityChanged.Notify(func() v1.AudioQuality {
		return playbackInfo.AudioQuality
	})

	// Free up resources taken up by previous stream
	playbin.SetState(gst.StateNull)
	playbin.SetArg("uri", "")

	OnStateChanged.Notify(func(state *State) {
		state.Status = StatusBuffering
		state.Position = 0
	})

	switch playbackInfo.ManifestMimeType {
	case v1.ManifestMimeTypeAudioMPD:
		enqueueMPDStream(playbackInfo)
	case v1.ManifestMimeTypeAudioBTS:
		enqueueBTSStream(playbackInfo)
	default:
		logger.Error("unsupported manifest mime type", "mime_type", playbackInfo.ManifestMimeType)
		return fmt.Errorf("unsupported manifest mime type: %s", playbackInfo.ManifestMimeType)
	}

	OnStateChanged.Notify(func(state *State) {
		state.Status = StatusPlaying
	})

	playbin.SetState(gst.StatePlaying)
	return nil
}
