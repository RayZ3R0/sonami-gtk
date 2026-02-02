package scrobbling

import (
	"log/slog"
	"time"

	"codeberg.org/dergs/tonearm/internal/player"
	"codeberg.org/dergs/tonearm/internal/signals"
)

var logger = slog.With("module", "scrobbler")

var TrackStarted = signals.NewStatelessSignal[*player.Track]()

var Scrobble = signals.NewStatelessSignal[*ScrobbleEvent]()

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
		go TrackStarted.Notify(t)

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
