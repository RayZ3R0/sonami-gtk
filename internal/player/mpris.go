package player

import (
	"fmt"

	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/pkg/mpris"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
	"github.com/godbus/dbus/v5"
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

	TrackChanged.Signal.On(func(trackInfo tonearm.Track) bool {
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

		artistNames := []string{}
		for _, artist := range trackInfo.Artists() {
			artistNames = append(artistNames, artist.Name())
		}

		// albumArtistsPaginator, err := album.Artists()
		// if err != nil {
		// 	return signals.Continue
		// }

		// albumArtists, err := albumArtistsPaginator.GetAll()
		// if err != nil {
		// 	return signals.Continue
		// }

		// albumArtistNames := []string{}
		// for _, artist := range albumArtists {
		// 	albumArtistNames = append(albumArtistNames, artist.Name())
		// }

		cover := trackInfo.Album().Cover(-1)
		if err != nil {
			return signals.Continue
		}

		mprisServer.SetTrackMetadata(map[string]any{
			"mpris:trackid": dbus.ObjectPath("/org/mpris/MediaPlayer2/TrackList/Track" + trackInfo.ID()),
			"mpris:artUrl":  cover,
			"mpris:length":  trackInfo.Duration().Microseconds(),
			"xesam:album":   trackInfo.Album().Title(),
			// "xesam:albumArtist": albumArtistNames,
			"xesam:artist": artistNames,
			"xesam:title":  trackInfo.Title(),
			"xesam:url":    fmt.Sprintf("https://tidal.com/track/%s", trackInfo.ID()),
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
