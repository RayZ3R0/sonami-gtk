package player

import (
	"sync"

	"codeberg.org/dergs/tidalwave/internal/signals"
)

var history = &History{
	Mutex:   sync.Mutex{},
	Current: nil,
	Entries: []*HistoryEntry{},
}

func init() {
	OnTrackChanged.On(func(trackInfo TrackInformation) bool {
		if trackInfo.ID != "" && (history.Current == nil || history.Current.TrackID != trackInfo.ID) {
			history.Push(&HistoryEntry{
				TrackID: trackInfo.ID,
			})
		}
		return signals.Continue
	})
}

type HistoryEntry struct {
	TrackID string
}

type History struct {
	sync.Mutex
	Current *HistoryEntry
	Entries []*HistoryEntry
}

func (h *History) Pop() *HistoryEntry {
	h.Lock()
	defer h.Unlock()

	if len(h.Entries) == 0 {
		return nil
	}

	h.Current = h.Entries[len(h.Entries)-1]
	h.Entries = h.Entries[:len(h.Entries)-1]

	return h.Current
}

func (h *History) Push(entry *HistoryEntry) {
	h.Lock()
	defer h.Unlock()

	if len(h.Entries) == 10 {
		h.Entries = h.Entries[1:]
	}

	if h.Current != nil {
		h.Entries = append(h.Entries, h.Current)
	}
	h.Current = entry
}
