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

type Codec string

const (
	CodecFLAC Codec = "Free Lossless Audio Codec (FLAC)"
	CodecAAC  Codec = "MPEG-4 AAC"
)

func (c Codec) String() string {
	switch c {
	case CodecFLAC:
		return "FLAC"
	case CodecAAC:
		return "AAC"
	default:
		return "unknown codec"
	}
}

type StreamQuality struct {
	Codec    Codec
	BitDepth int

	// BitRate is the bitrate in bps if the codec is AAC, otherwise 0.
	BitRate uint32

	// SampleRate is the sample rate in Hz if the codec is FLAC, otherwise 0.
	SampleRate int
}
