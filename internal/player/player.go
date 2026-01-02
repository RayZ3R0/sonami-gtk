package player

import (
	"context"
	"log/slog"
	"math/rand/v2"

	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"github.com/go-gst/go-gst/gst"
	"github.com/infinytum/injector"
)

var (
	playbin *gst.Element
	logger  *slog.Logger
)

func init() {
	gst.Init(nil)
	logger = slog.With("module", "player")
	pb, err := gst.NewElement("playbin")
	if err != nil {
		panic(err)
	}
	playbin = pb

	playbin.GetBus().AddWatch(onBusMessage)
	playbin.Connect("notify::volume", onVolumeChange)
	onVolumeChange()
}

func PlayTrack(trackId string) error {
	tidal, err := injector.Inject[*tidalapi.TidalAPI]()
	if err != nil {
		return err
	}

	track, err := tidal.OpenAPI.V2.Tracks.Track(context.Background(), trackId, "albums.coverArt", "artists")
	if err != nil {
		return err
	}

	BaseQueue.Clear()

	return playTrack(track)
}

func PlayPlaylist(playlistId string, shuffle bool) error {
	tidal, err := injector.Inject[*tidalapi.TidalAPI]()
	if err != nil {
		return err
	}

	items, err := tidal.OpenAPI.V2.Playlists.Items(context.Background(), playlistId, "", "items", "items.artists", "items.albums.coverArt")
	if err != nil {
		return err
	}

	BaseQueue.Clear()

	tracks := items.Included.Tracks(items.Data...)

	if shuffle {
		rand.Shuffle(len(tracks), func(i, j int) {
			tracks[i], tracks[j] = tracks[j], tracks[i]
		})
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
