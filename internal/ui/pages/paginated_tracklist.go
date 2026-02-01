package pages

import (
	"log/slog"
	"unsafe"

	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/internal/ui/components/tracklist"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/pagination"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func NewPaginatedTracklistPage(
	paginator *pagination.Paginator[openapi.Track],
	factory func() *tracklist.TrackList[*openapi.Track],
	styleFactory func(*tracklist.TrackList[*openapi.Track]) schwifty.BaseWidgetable,
) (schwifty.ScrolledWindow, error) {
	firstPage, err := paginator.GetFirstPage()
	if err != nil {
		return nil, err
	}

	list := factory()
	for _, track := range firstPage {
		list.AddTrack(&track)
	}

	return ScrolledWindow().
		Child(styleFactory(list)).
		Policy(gtk.PolicyNeverValue, gtk.PolicyAutomaticValue).
		ConnectReachEdgeSoon(gtk.PosBottomValue, func() bool {
			if !paginator.IsConsumed() {
				items, err := paginator.Next()
				if err != nil {
					return signals.Continue
				}

				schwifty.OnMainThreadOnce(func(u uintptr) {
					var list *tracklist.TrackList[*openapi.Track]
					list = (*tracklist.TrackList[*openapi.Track])(unsafe.Pointer(u))
					for _, track := range items {
						list.AddTrack(&track)
					}
				}, uintptr(unsafe.Pointer(list)))
			} else {
				slog.Debug("No more items to fetch")
				return signals.Unsubscribe
			}
			return signals.Continue
		}), nil
}
