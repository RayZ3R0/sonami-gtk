package player

import (
	"context"
	"fmt"

	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	v1 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v1"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/v1/tracks"
	"github.com/go-gst/go-gst/gst"
	"github.com/infinytum/injector"
)

var (
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
	onVolumeChange()
}

func Play(trackId int) error {
	tidal, err := injector.Inject[*tidalapi.TidalAPI]()
	if err != nil {
		return err
	}

	track, err := tidal.V1.Tracks.Track(context.Background(), trackId)
	if err != nil {
		return err
	}

	playbackInfo, err := tidal.V1.Tracks.PlaybackInfo(context.Background(), trackId, tracks.PlaybackInfoOptions{})
	if err != nil {
		return err
	}

	playbin.SetState(gst.StateReady)
	OnState.Notify(func(state *State) {
		state.Track = track
		state.Status = StatusBuffering
		state.Duration = track.Duration
		state.Position = 0
	})

	switch playbackInfo.ManifestMimeType {
	case v1.ManifestMimeTypeAudioMPD:
		enqueueMPDStream(playbackInfo)
	default:
		return fmt.Errorf("unsupported manifest mime type: %s", playbackInfo.ManifestMimeType)
	}

	OnState.Notify(func(state *State) {
		state.Status = StatusPlaying
	})

	playbin.SetState(gst.StatePlaying)
	return nil
}
