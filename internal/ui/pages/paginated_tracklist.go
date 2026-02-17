package pages

import (
	"log/slog"
	"unsafe"

	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/internal/ui/components"
	"codeberg.org/dergs/tonearm/internal/ui/components/tracklist"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func NewPaginatedTracklistPage(
	paginator tonearm.Paginator[tonearm.Track],
	factory func() *tracklist.TrackList,
	styleFactory func(*tracklist.TrackList) schwifty.BaseWidgetable,
) (schwifty.ScrolledWindow, error) {
	firstPage, err := paginator.NextPage()
	if err != nil {
		return nil, err
	}

	list := factory()
	for _, track := range firstPage {
		list.AddTrack(track)
	}

	return ScrolledWindow().
		Child(
			components.MainContent(
				styleFactory(list),
			),
		).
		Policy(gtk.PolicyNeverValue, gtk.PolicyAutomaticValue).
		ConnectReachEdgeSoon(gtk.PosBottomValue, func() bool {
			if !paginator.IsConsumed() {
				items, err := paginator.NextPage()
				if err != nil {
					return signals.Continue
				}

				schwifty.OnMainThreadOnce(func(u uintptr) {
					var list *tracklist.TrackList
					list = (*tracklist.TrackList)(unsafe.Pointer(u))
					for _, track := range items {
						list.AddTrack(track)
					}
				}, uintptr(unsafe.Pointer(list)))
			} else {
				slog.Debug("No more items to fetch")
				return signals.Unsubscribe
			}
			return signals.Continue
		}), nil
}
