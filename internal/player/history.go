package player

import (
	"sync"

	"codeberg.org/dergs/tidalwave/internal/signals"
)

var history = &History{
	Mutex:   sync.Mutex{},
	Current: signals.NewStatefulSignal[*HistoryEntry](nil),
	Entries: signals.NewStatefulSignal[[]*HistoryEntry](nil),
}

type HistoryEntry struct {
	TrackID string
}

type History struct {
	sync.Mutex
	Current *signals.StatefulSignal[*HistoryEntry]
	Entries *signals.StatefulSignal[[]*HistoryEntry]
}

func (h *History) Pop() *HistoryEntry {
	h.Lock()
	defer h.Unlock()

	if len(h.Entries.CurrentValue()) == 0 {
		return nil
	}

	h.Entries.Notify(func(entries []*HistoryEntry) []*HistoryEntry {
		h.Current.Notify(func(oldValue *HistoryEntry) *HistoryEntry {
			return entries[len(entries)-1]
		})
		return entries[:len(entries)-1]
	})

	return h.Current.CurrentValue()
}

func (h *History) Push(entry *HistoryEntry) {
	h.Lock()
	defer h.Unlock()

	if len(h.Entries.CurrentValue()) == 10 {
		h.Entries.Notify(func(entries []*HistoryEntry) []*HistoryEntry {
			return entries[1:]
		})
	}

	if current := h.Current.CurrentValue(); current != nil {
		h.Entries.Notify(func(entries []*HistoryEntry) []*HistoryEntry {
			return append(entries, current)
		})
	}

	h.Current.Notify(func(oldValue *HistoryEntry) *HistoryEntry {
		return entry
	})
}
