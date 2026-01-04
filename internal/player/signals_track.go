package player

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"codeberg.org/dergs/tidalwave/internal/signals"
	"codeberg.org/dergs/tidalwave/pkg/mpris"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	"github.com/infinytum/injector"
)

func init() {
	OnTrackChanged.Signal.On(func(trackInfo TrackInformation) bool {
		mpris := injector.MustInject[*mpris.Server]()
		mpris.EnableControl()

		artists := []string{}
		for _, artist := range trackInfo.Artists {
			artists = append(artists, artist.Attributes.Name)
		}

		album := trackInfo.Albums[0]
		albumArtists := []string{}
		for _, artist := range album.Included.Artists(album.Data.Relationships.Artists.Data...) {
			albumArtists = append(albumArtists, artist.Data.Attributes.Name)
		}

		mpris.SetTrackMetadata(map[string]any{
			"mpris:artUrl":      trackInfo.CoverURL,
			"mpris:length":      trackInfo.Duration.Microseconds(),
			"xesam:album":       album.Data.Attributes.Title,
			"xesam:albumArtist": albumArtists,
			"xesam:artist":      artists,
			"xesam:title":       trackInfo.Title,
			"xesam:url":         fmt.Sprintf("https://tidal.com/track/%s", trackInfo.ID),
		})

		return signals.Continue
	})
}

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
	Albums   []openapi.Album
	CoverURL string
	Duration time.Duration
	ID       string
	Title    string
}

func (t TrackInformation) ArtistNames() string {
	names := make([]string, len(t.Artists))
	for i, artist := range t.Artists {
		names[i] = artist.Attributes.Name
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
