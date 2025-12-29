package router

import (
	"maps"
	"sync"

	"github.com/jwijenbergh/puregotk/v4/gtk"
)

type HistoryEntry struct {
	Path      string
	Params    Params
	Response  *Response
	PageTitle string
	View      *gtk.Widget
	Toolbar   *gtk.Widget
}

type History struct {
	sync.Mutex
	array     []HistoryEntry
	maxLength int
}

func (h *History) Current() *HistoryEntry {
	h.Lock()
	defer h.Unlock()

	if len(h.array) < 1 {
		return nil
	}

	return &h.array[len(h.array)-1]
}

func (h *History) IsCurrentlyOn(path string, params Params) bool {
	if len(h.array) < 1 {
		return false
	}

	lastItem := h.array[len(h.array)-1]
	if lastItem.Path != path {
		return false
	}

	return maps.Equal(lastItem.Params, params)
}

func (h *History) Length() int {
	return len(h.array)
}

func (h *History) Push(entry HistoryEntry) {
	h.Lock()
	defer h.Unlock()
	defer HistoryUpdated.Notify(h)

	h.array = append(h.array, entry)
	if len(h.array) > h.maxLength {
		h.array = h.array[len(h.array)-h.maxLength:]
	}
}

func (h *History) Pop() *HistoryEntry {
	h.Lock()
	defer h.Unlock()
	defer HistoryUpdated.Notify(h)

	if len(h.array) < 1 {
		return nil
	}

	h.array = h.array[:len(h.array)-1]
	previous := &h.array[len(h.array)-1]

	return previous
}

var history = &History{
	sync.Mutex{},
	make([]HistoryEntry, 0),
	10,
}
