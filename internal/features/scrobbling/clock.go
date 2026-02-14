package scrobbling

import (
	"context"
	"time"

	"codeberg.org/dergs/tonearm/internal/player"
)

type Clock struct {
	counter   int
	track     *player.Track
	isRunning bool
	startedAt time.Time

	context    context.Context
	cancelFunc context.CancelFunc
}

func (c *Clock) Start() {
	if c.isRunning {
		return
	}
	if c.counter >= int(c.track.Duration.Seconds())/2 {
		return
	}
	c.isRunning = true

	ctx, cancel := context.WithCancel(context.Background())
	c.context = ctx
	c.cancelFunc = cancel

	go func() {
		defer func() {
			c.isRunning = false
		}()

		for {
			select {
			case <-c.context.Done():
				return
			case <-time.After(time.Second):
				c.counter++
				if c.counter >= int(c.track.Duration.Seconds())/2 || c.counter >= int((time.Minute*4).Seconds()) {
					logger.Debug("notifying scrobblers that a track should be scrobbled")

					event := new(ScrobbleEvent)
					event.Track = c.track
					event.ListenedAt = c.startedAt

					for _, scrobbler := range Scrobblers {
						if !scrobbler.IsConfigured() {
							continue
						}

						go scrobbler.Scrobble(event)
					}

					c.cancelFunc()
					return
				}
			}
		}
	}()
}

func (c *Clock) Stop() {
	if !c.isRunning {
		return
	}
	c.cancelFunc()
}

func newClock(track *player.Track) *Clock {
	return &Clock{
		track:     track,
		startedAt: time.Now(),
	}
}
