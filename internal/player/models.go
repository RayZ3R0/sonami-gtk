package player

import (
	"fmt"
	"time"

	"codeberg.org/dergs/tonearm/internal/router"
	"codeberg.org/dergs/tonearm/internal/ui/sidebar/navigation"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
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

var UserQueueSource = userQueueSource{}

type userQueueSource struct{}

func (u userQueueSource) Cover(preferredSize int) string {
	return ""
}

func (u userQueueSource) Route() string {
	return fmt.Sprintf("%s%s", router.SidebarNavPrefix, navigation.PathQueue)
}

func (u userQueueSource) Title() string {
	return "Your Queue"
}

func (u userQueueSource) SourceType() tonearm.SourceType {
	return tonearm.SourceTypeQueue
}
