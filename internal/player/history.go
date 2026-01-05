package player

import (
	"sync"

	"codeberg.org/dergs/tidalwave/internal/signals"
)

var (
	unmanagedHistory = &History{
		Mutex:   sync.Mutex{},
		Signal:  signals.NewSignal[func(list []*HistoryEntry, current *HistoryEntry) bool](),
		Current: nil,
		Entries: []*HistoryEntry{},
	}
	managedHistory = &History{
		Mutex:   sync.Mutex{},
		Signal:  signals.NewSignal[func(list []*HistoryEntry, current *HistoryEntry) bool](),
		Current: nil,
		Entries: []*HistoryEntry{},
	}
)

func init() {
	OnTrackChanged.On(func(trackInfo TrackInformation) bool {
		history := unmanagedHistory
		if currentHistoryType == HistoryTypeManaged {
			history = managedHistory
		}

		if trackInfo.ID != "" && (history.Current == nil || history.Current.TrackID != trackInfo.ID) {
			history.Push(&HistoryEntry{
				TrackID: trackInfo.ID,
			})
		}
		return signals.Continue
	})
}

type HistorySignal struct {
	Signal  *signals.Signal[func(list []*HistoryEntry, current *HistoryEntry) bool]
	history *History
}

var (
	OnUnmanagedHistoryChanged = &HistorySignal{
		Signal:  &unmanagedHistory.Signal,
		history: unmanagedHistory,
	}
	OnManagedHistoryChanged = &HistorySignal{
		Signal:  &managedHistory.Signal,
		history: managedHistory,
	}
)

type HistoryType int

const (
	HistoryTypeManaged HistoryType = iota
	HistoryTypeUnmanaged
)

var currentHistoryType = HistoryTypeUnmanaged

type HistoryEntry struct {
	TrackID string
}

type History struct {
	sync.Mutex
	Signal  signals.Signal[func([]*HistoryEntry, *HistoryEntry) bool]
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

	h.Signal.Notify(h.Entries, h.Current)
	return h.Current
}

func (h *History) Push(entry *HistoryEntry) {
	h.Lock()
	defer h.Unlock()

	if h.Current != nil {
		h.Entries = append(h.Entries, h.Current)
	}
	h.Current = entry
	h.Signal.Notify(h.Entries, h.Current)
}

func (h *History) Clear() {
	h.Lock()
	defer h.Unlock()

	h.Entries = []*HistoryEntry{}
	h.Current = nil

	h.Signal.Notify(h.Entries, h.Current)
}

func (h *HistorySignal) On(callback func(list []*HistoryEntry, current *HistoryEntry) bool) {
	callback(h.history.Entries, h.history.Current)
	h.Signal.On(callback)
}
