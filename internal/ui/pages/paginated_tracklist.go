package pages

import (
	"log/slog"
	"unsafe"

	"github.com/RayZ3R0/sonami-gtk/internal/signals"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components/tracklist"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
	"codeberg.org/puregotk/puregotk/v4/gtk"
)

func NewPaginatedTracklistPage(
	paginator sonami.Paginator[sonami.Track],
	styleFactory func(*tracklist.TrackList) schwifty.BaseWidgetable,
	columns ...tracklist.ColumnFunc,
) (schwifty.ScrolledWindow, error) {
	firstPage, err := paginator.NextPage()
	if err != nil {
		return nil, err
	}

	list := tracklist.NewTrackList(columns...)
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
