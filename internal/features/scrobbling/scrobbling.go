package scrobbling

import (
	"log/slog"
	"time"

	"codeberg.org/dergs/tonearm/internal/player"
	"codeberg.org/dergs/tonearm/internal/signals"
)

var logger = slog.With("module", "scrobbler")

var Scrobblers []Scrobbler

type Scrobbler interface {
	NowPlaying(*player.Track)
	Scrobble(*ScrobbleEvent)
	IsConfigured() bool
}

type ScrobbleEvent struct {
	Track      *player.Track
	ListenedAt time.Time
}

func init() {
	var scrobbleClock *Clock
	player.TrackChanged.On(func(t *player.Track) bool {
		if t == nil {
			return signals.Continue
		}

		logger.Debug("notifying scrobblers that a new track has started playing")
		for _, scrobbler := range Scrobblers {
			if !scrobbler.IsConfigured() {
				continue
			}

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
