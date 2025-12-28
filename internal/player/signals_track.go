package player

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"codeberg.org/dergs/tidalwave/internal/signals"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
)

var OnTrackChanged = trackChangedSignal{
	signals.NewSignal[func(trackInfo TrackInformation) bool](),
	TrackInformation{},
	sync.Mutex{},
}

type trackChangedSignal struct {
	signals.Signal[func(trackInfo TrackInformation) bool]
	current TrackInformation
	lock    sync.Mutex
}

func (r *trackChangedSignal) Notify(callback func(trackInfo *TrackInformation)) {
	r.lock.Lock()
	defer r.lock.Unlock()
	newState := r.current
	callback(&newState)
	if newState.Equals(r.current) {
		return
	}
	logger.Info("track changed", "track_id", newState.ID, "track_information", newState)
	r.current = newState
	r.Signal.Notify(newState)
}

func (r *trackChangedSignal) On(handler func(trackInfo TrackInformation) bool) *signals.Subscription {
	handler(r.current)
	return r.Signal.On(handler)
}

type TrackInformation struct {
	Artists  []openapi.ArtistAttributes
	CoverURL string
	Duration time.Duration
	ID       int
	Title    string
}

func (t TrackInformation) ArtistNames() string {
	names := make([]string, len(t.Artists))
	for i, artist := range t.Artists {
		names[i] = artist.Name
	}
	return strings.Join(names, ", ")
}

func (t TrackInformation) Equals(other TrackInformation) bool {
	return t.ID == other.ID
}

func (t TrackInformation) String() string {
	return fmt.Sprintf("%s by %s - %s", t.Title, t.ArtistNames(), t.Duration)
}

type TrackInformationArtist struct {
	Name string
	ID   int
}
