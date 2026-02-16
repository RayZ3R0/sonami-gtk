package player

import (
	"time"
)

type RepeatMode int

const (
	RepeatModeNone RepeatMode = iota
	RepeatModeQueue
	RepeatModeTrack
)

type PlaybackStatus string

const (
	PlaybackStatusPlaying PlaybackStatus = "playing"
	PlaybackStatusPaused  PlaybackStatus = "paused"
	PlaybackStatusStopped PlaybackStatus = "stopped"
)

type PlaybackState struct {
	// Expected duration of the currently playing stream as reported by playbin
	Duration time.Duration

	// Whether the position was changed by the user
	IsSeeking bool

	// Current position of the playback as reported by playbin
	Position time.Duration

	Status PlaybackStatus

	Loading bool
}
