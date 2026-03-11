package discord

import (
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/RayZ3R0/sonami-gtk/internal/player"
	"github.com/RayZ3R0/sonami-gtk/internal/settings"
	"github.com/RayZ3R0/sonami-gtk/internal/signals"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
)

const (
	appID           = "1459143320604508251"
	seekThreshold   = 4 * time.Second
	pauseClearDelay = 30 * time.Second
)

var logger = slog.With("module", "discord-rpc")

type snapshot struct {
	track  sonami.Track
	status player.PlaybackStatus
	pos    time.Duration
}

var (
	mu      sync.Mutex
	current snapshot
	changed = make(chan struct{}, 1)
)

func notify() {
	select {
	case changed <- struct{}{}:
	default:
	}
}

func init() {
	// Initialize settings on the main goroutine (init() runs from main goroutine)
	// to avoid g.Lazy data races.
	discordSettings := settings.Discord()

	go run(discordSettings)

	player.TrackChanged.On(func(t sonami.Track) bool {
		mu.Lock()
		current.track = t
		current.pos = 0
		current.status = player.PlaybackStatusStopped
		mu.Unlock()
		notify()
		return signals.Continue
	})

	player.PlaybackStateChanged.On(func(ps *player.PlaybackState) bool {
		if ps == nil {
			return signals.Continue
		}
		mu.Lock()
		statusChanged := current.status != ps.Status
		current.status = ps.Status
		current.pos = ps.Position
		mu.Unlock()
		if statusChanged {
			notify()
		}
		return signals.Continue
	})
}

func run(s *settings.DiscordSettings) {
	var (
		client     *ipcClient
		lastSent   *activity
		pid        = os.Getpid()
		backoff    = time.Second
		clearTimer = (<-chan time.Time)(nil) // fires 30s after pause
	)

	disconnect := func() {
		if client != nil {
			client.close()
			client = nil
		}
		lastSent = nil
		clearTimer = nil
	}

	for {
		// Wait for a state change, 30s pause-clear timer, or a 10s background poll.
		select {
		case <-changed:
		case <-clearTimer:
			clearTimer = nil
			if client != nil {
				if err := client.setActivity(pid, nil); err != nil {
					logger.Debug("discord: clear on pause timeout failed", "error", err)
					disconnect()
				} else {
					lastSent = nil
				}
			}
			continue
		case <-time.After(10 * time.Second):
		}

		if !s.RichPresenceEnabled() {
			if client != nil {
				_ = client.setActivity(pid, nil)
				disconnect()
			}
			lastSent = nil
			clearTimer = nil
			continue
		}

		// Connect if needed.
		if client == nil {
			var err error
			client, err = dialIPC()
			if err != nil {
				logger.Debug("discord: IPC unavailable", "error", err)
				time.Sleep(backoff)
				if backoff < 30*time.Second {
					backoff *= 2
				}
				continue
			}
			if err := client.handshake(appID); err != nil {
				logger.Debug("discord: handshake failed", "error", err)
				disconnect()
				time.Sleep(backoff)
				if backoff < 30*time.Second {
					backoff *= 2
				}
				continue
			}
			backoff = time.Second
			lastSent = nil
		}

		mu.Lock()
		snap := current
		mu.Unlock()

		a := buildActivity(snap, lastSent)
		if activityEqual(a, lastSent) {
			continue
		}

		if err := client.setActivity(pid, a); err != nil {
			logger.Debug("discord: setActivity failed", "error", err)
			disconnect()
			continue
		}
		lastSent = a

		// Start the 30-second pause-clear timer when paused, cancel it otherwise.
		if a != nil && snap.status == player.PlaybackStatusPaused {
			clearTimer = time.After(pauseClearDelay)
		} else {
			clearTimer = nil
		}
	}
}

// buildActivity constructs the Discord activity for the given snapshot.
// lastSent is used for seek detection — returns the same pointer when
// only the start timestamp needs recalibration due to an observed seek.
func buildActivity(snap snapshot, lastSent *activity) *activity {
	if snap.track == nil || snap.status == player.PlaybackStatusStopped {
		return nil
	}

	artists := strings.Join(snap.track.Artists().Names(), ", ")
	if artists == "" {
		artists = "Unknown Artist"
	}
	title := sonami.FormatTitle(snap.track)
	album := snap.track.Album().Title()
	coverURL := snap.track.Cover(512)

	a := &activity{
		Type:    2, // Listening
		Details: title,
		Assets: &assets{
			LargeImage: coverURL,
			LargeText:  album,
		},
	}

	switch snap.status {
	case player.PlaybackStatusPlaying:
		now := time.Now().Unix()
		start := now - int64(snap.pos.Seconds())
		dur := snap.track.Duration()
		a.State = artists
		if dur > 0 {
			a.Timestamps = &timestamps{
				Start: start,
				End:   start + int64(dur.Seconds()),
			}
		}
	case player.PlaybackStatusPaused:
		a.State = artists + " • Paused"
	}

	return a
}

// activityEqual returns true when a and b would produce the same Discord display.
// For playing state the start timestamp is coarsened to seekThreshold buckets so
// that normal playback drift does not trigger unnecessary resyncs.
func activityEqual(a, b *activity) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	if a == nil {
		return true
	}
	if a.Details != b.Details || a.State != b.State {
		return false
	}
	if (a.Timestamps == nil) != (b.Timestamps == nil) {
		return false
	}
	if a.Timestamps != nil && b.Timestamps != nil {
		bucketSecs := int64(seekThreshold.Seconds())
		if a.Timestamps.Start/bucketSecs != b.Timestamps.Start/bucketSecs {
			return false
		}
	}
	return true
}
