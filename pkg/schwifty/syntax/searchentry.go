package syntax

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/callback"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func SearchEntry() schwifty.SearchEntry {
	return managed("SearchEntry", func() *gtk.SearchEntry {
		searchEntry := gtk.NewSearchEntry()
		searchEntry.ConnectActivate(&callback.SearchEntryActivateCallback)
		searchEntry.ConnectSearchChanged(&callback.SearchChangedCallback)
		return searchEntry
	})
}
