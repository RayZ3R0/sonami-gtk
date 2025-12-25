package player

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	v1 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v1"
	tracksv1 "codeberg.org/dergs/tidalwave/pkg/tidalapi/v1/tracks"
	"github.com/go-gst/go-gst/gst"
	"github.com/infinytum/injector"
)

var (
	playbin *gst.Element
	logger  *slog.Logger
)

func init() {
	gst.Init(nil)
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})).With("module", "player")
	pb, err := gst.NewElement("playbin")
	if err != nil {
		panic(err)
	}
	playbin = pb

	playbin.GetBus().AddWatch(onBusMessage)
	playbin.Connect("notify::volume", onVolumeChange)
	onVolumeChange()
}

func Play(trackId int) error {
	tidal, err := injector.Inject[*tidalapi.TidalAPI]()
	if err != nil {
		return err
	}

	openTrack, err := tidal.OpenAPI.V2.Tracks.Track(context.Background(), trackId, "albums.coverArt", "artists")
	if err != nil {
		return err
	}
	OnTrackChanged.Notify(func(trackInfo *TrackInformation) {
		trackInfo.Artists = []openapi.ArtistAttributes{}
		trackInfo.CoverURL = ""
		trackInfo.Duration = openTrack.Data.Attributes.Duration.Duration
		trackInfo.ID = trackId
		trackInfo.Title = openTrack.Data.Attributes.Title

		for _, artist := range openTrack.Included.PlainArtists(openTrack.Data.Relationships.Artists.Data...) {
			trackInfo.Artists = append(trackInfo.Artists, artist.Attributes)
		}

		for _, album := range openTrack.Included.Albums(openTrack.Data.Relationships.Albums.Data...) {
			for _, artwork := range album.Included.PlainArtworks(album.Data.Relationships.CoverArt.Data...) {
				trackInfo.CoverURL = artwork.Attributes.Files.AtLeast(320).Href
			}
		}
	})

	playbackInfo, err := tidal.V1.Tracks.PlaybackInfo(context.Background(), trackId, tracksv1.PlaybackInfoOptions{})
	if err != nil {
		return err
	}

	playbin.SetState(gst.StateReady)
	OnStateChanged.Notify(func(state *State) {
		state.Status = StatusBuffering
		state.Position = 0
	})

	switch playbackInfo.ManifestMimeType {
	case v1.ManifestMimeTypeAudioMPD:
		enqueueMPDStream(playbackInfo)
	default:
		return fmt.Errorf("unsupported manifest mime type: %s", playbackInfo.ManifestMimeType)
	}

	OnStateChanged.Notify(func(state *State) {
		state.Status = StatusPlaying
	})

	playbin.SetState(gst.StatePlaying)
	return nil
}
