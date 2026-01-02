package player

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"codeberg.org/dergs/tidalwave/internal/signals"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	v1 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v1"
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
	Artists  []openapi.ArtistData
	CoverURL string
	Duration time.Duration
	ID       string
	Title    string
	Quality  v1.AudioQuality
}

func (t TrackInformation) ArtistNames() string {
	names := make([]string, len(t.Artists))
	for i, artist := range t.Artists {
		names[i] = artist.Attributes.Name
	}
	return strings.Join(names, ", ")
}

func (t TrackInformation) Equals(other TrackInformation) bool {
	return t.ID == other.ID && t.Quality == other.Quality
}

func (t TrackInformation) String() string {
	return fmt.Sprintf("%s by %s - %s", t.Title, t.ArtistNames(), t.Duration)
}

type TrackInformationArtist struct {
	Name string
	ID   int
}
