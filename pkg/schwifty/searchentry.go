package schwifty

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty/callback"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

//go:generate go run codeberg.org/dergs/tidalwave/pkg/schwifty/gen SearchEntry *gtk.SearchEntry

func (f SearchEntry) ConnectActivate(cb func(gtk.SearchEntry)) SearchEntry {
	return func() *gtk.SearchEntry {
		searchEntry := f()
		callback.HandleCallback(searchEntry.Widget, "activate", cb)
		return searchEntry
	}
}

func (f SearchEntry) ConnectSearchChanged(cb func(gtk.SearchEntry)) SearchEntry {
	return func() *gtk.SearchEntry {
		searchEntry := f()
		callback.HandleCallback(searchEntry.Widget, "search-changed", cb)
		return searchEntry
	}
}

func (f SearchEntry) PlaceholderText(text string) SearchEntry {
	return func() *gtk.SearchEntry {
		searchEntry := f()
		searchEntry.SetPlaceholderText(text)
		return searchEntry
	}
}

func (f SearchEntry) SearchDelay(delay uint) SearchEntry {
	return func() *gtk.SearchEntry {
		searchEntry := f()
		searchEntry.SetSearchDelay(delay)
		return searchEntry
	}
}
