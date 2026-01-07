package player

import (
	"fmt"

	"codeberg.org/dergs/tidalwave/internal/signals"
	"codeberg.org/dergs/tidalwave/pkg/mpris"
	"github.com/infinytum/injector"
)

func init() {
	PlaybackStateChanged.Signal.On(func(state *PlaybackState) bool {
		mprisClient, err := injector.Inject[*mpris.Server]()
		if err != nil {
			return signals.Continue
		}

		mprisClient.SetPosition(state.Position)
		switch state.Status {
		case PlaybackStatusBuffering, PlaybackStatusPaused:
			mprisClient.SetPlaybackStatus(mpris.PlaybackStatusPaused)
		case PlaybackStatusPlaying:
			mprisClient.SetPlaybackStatus(mpris.PlaybackStatusPlaying)
		case PlaybackStatusStopped:
			mprisClient.SetPlaybackStatus(mpris.PlaybackStatusStopped)
		}

		return signals.Continue
	})
	TrackChanged.Signal.On(func(trackInfo *Track) bool {
		mpris, err := injector.Inject[*mpris.Server]()
		if err != nil {
			return signals.Continue
		}

		mpris.EnableControl()
		if trackInfo == nil {
			return signals.Continue
		}

		artists := []string{}
		for _, artist := range trackInfo.Artists {
			artists = append(artists, artist.Attributes.Name)
		}

		album := trackInfo.Albums[0]
		albumArtists := []string{}
		for _, artist := range album.Included.Artists(album.Data.Relationships.Artists.Data...) {
			albumArtists = append(albumArtists, artist.Data.Attributes.Name)
		}

		mpris.SetTrackMetadata(map[string]any{
			"mpris:artUrl":      trackInfo.CoverURL,
			"mpris:length":      trackInfo.Duration.Microseconds(),
			"xesam:album":       album.Data.Attributes.Title,
			"xesam:albumArtist": albumArtists,
			"xesam:artist":      artists,
			"xesam:title":       trackInfo.Title,
			"xesam:url":         fmt.Sprintf("https://tidal.com/track/%s", trackInfo.ID),
		})

		return signals.Continue
	})
	VolumeChanged.Signal.On(func(volume float64) bool {
		server, err := injector.Inject[*mpris.Server]()
		if err != nil {
			return signals.Continue
		}

		server.SetVolume(volume)
		return signals.Continue
	})
}
