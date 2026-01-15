package player

import (
	"fmt"
	"strings"
	"time"

	"codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"
)

type RepeatMode int

const (
	RepeatModeNone RepeatMode = iota
	RepeatModeQueue
	RepeatModeTrack
)

type PlaybackStatus string

const (
	PlaybackStatusLoadingTrack PlaybackStatus = "loading"
	PlaybackStatusBuffering    PlaybackStatus = "buffering"
	PlaybackStatusPlaying      PlaybackStatus = "playing"
	PlaybackStatusPaused       PlaybackStatus = "paused"
	PlaybackStatusStopped      PlaybackStatus = "stopped"
)

type PlaybackState struct {
	// Expected duration of the currently playing stream as reported by playbin
	Duration time.Duration

	// Whether the position was changed by the user
	IsSeeking bool

	// Current position of the playback as reported by playbin
	Position time.Duration

	Status PlaybackStatus
}

type Track struct {
	Artists  []openapi.ArtistData
	Albums   []openapi.Album
	CoverURL string
	Duration time.Duration
	ID       string
	ISRC     string
	Title    string
}

func (t Track) ArtistNames() string {
	names := make([]string, len(t.Artists))
	for i, artist := range t.Artists {
		names[i] = artist.Attributes.Name
	}
	return strings.Join(names, ", ")
}

func (t Track) String() string {
	return fmt.Sprintf("%s by %s - %s", t.Title, t.ArtistNames(), t.Duration)
}
