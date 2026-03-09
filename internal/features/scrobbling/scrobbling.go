package scrobbling

import (
	"log/slog"
	"time"

	"github.com/RayZ3R0/sonami-gtk/internal/player"
	"github.com/RayZ3R0/sonami-gtk/internal/signals"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
)

var logger = slog.With("module", "scrobbler")

var Scrobblers []Scrobbler

type Scrobbler interface {
	NowPlaying(sonami.Track)
	Scrobble(*ScrobbleEvent)
	IsConfigured() bool
	GetName() string
}

type ScrobbleEvent struct {
	Track      sonami.Track
	ListenedAt time.Time
}

func init() {
	var scrobbleClock *Clock
	player.TrackChanged.On(func(t sonami.Track) bool {
		if t == nil {
			return signals.Continue
		}

		logger.Debug("notifying scrobblers that a new track has started playing")
		for _, scrobbler := range Scrobblers {
			if !scrobbler.IsConfigured() {
				logger.Debug("skipping now playing event", "service", scrobbler.GetName())
				continue
			}

			logger.Debug("sending NowPlaying event", "service", scrobbler.GetName())
			go scrobbler.NowPlaying(t)
		}

		if scrobbleClock != nil {
			scrobbleClock.Stop()
		}
		scrobbleClock = newClock(t)
		return signals.Continue
	})

	player.PlaybackStateChanged.On(func(ps *player.PlaybackState) bool {
		if scrobbleClock == nil {
			return signals.Continue
		}

		if ps.Status == player.PlaybackStatusPlaying {
			scrobbleClock.Start()
		} else {
			scrobbleClock.Stop()
		}
		return signals.Continue
	})
}
