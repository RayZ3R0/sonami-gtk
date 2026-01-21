package player

import (
	"fmt"

	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/pkg/mpris"
	"github.com/infinytum/injector"
)

func init() {
	PlaybackStateChanged.Signal.On(func(state *PlaybackState) bool {
		mprisServer, err := injector.Inject[*mpris.Server]()
		if err != nil {
			return signals.Continue
		}

		mprisServer.SetPosition(state.Position, state.IsSeeking)
		if state.Loading {
			mprisServer.SetPlaybackStatus(mpris.PlaybackStatusPaused)
		} else {
			switch state.Status {
			case PlaybackStatusPaused:
				mprisServer.SetPlaybackStatus(mpris.PlaybackStatusPaused)
			case PlaybackStatusPlaying:
				mprisServer.SetPlaybackStatus(mpris.PlaybackStatusPlaying)
			case PlaybackStatusStopped:
				mprisServer.SetPlaybackStatus(mpris.PlaybackStatusStopped)
			}
		}

		return signals.Continue
	})

	RepeatModeChanged.On(func(rm RepeatMode) bool {
		mprisServer, err := injector.Inject[*mpris.Server]()
		if err != nil {
			return signals.Continue
		}

		switch rm {
		case RepeatModeNone:
			mprisServer.SetLoopStatus(mpris.LoopNone)
		case RepeatModeTrack:
			mprisServer.SetLoopStatus(mpris.LoopTrack)
		case RepeatModeQueue:
			mprisServer.SetLoopStatus(mpris.LoopPlaylist)
		}
		return signals.Continue
	})

	TrackChanged.Signal.On(func(trackInfo *Track) bool {
		mprisServer, err := injector.Inject[*mpris.Server]()
		if err != nil {
			return signals.Continue
		}

		if trackInfo == nil {
			mprisServer.SetTrackMetadata(map[string]any{})
			mprisServer.Disconnect()

			return signals.Continue
		}

		mprisServer.Connect()

		artists := []string{}
		for _, artist := range trackInfo.Artists {
			artists = append(artists, artist.Attributes.Name)
		}

		album := trackInfo.Albums[0]
		albumArtists := []string{}
		for _, artist := range album.Included.Artists(album.Data.Relationships.Artists.Data...) {
			albumArtists = append(albumArtists, artist.Data.Attributes.Name)
		}

		mprisServer.SetTrackMetadata(map[string]any{
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
		mprisServer, err := injector.Inject[*mpris.Server]()
		if err != nil {
			return signals.Continue
		}

		mprisServer.SetVolume(volume)
		return signals.Continue
	})
}
